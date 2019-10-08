#!/bin/bash

set -e

usage() {
  echo "Usage:
          -u  |  --url                 - Required. URL to wait for a 200 OK.
          -c  |  --cookie-jar          - Optional. Cookie jar to use to store cookies.
          -s  |  --successful-requests - Optional. Number of successful requests to wait for. Default 3.
          -t  |  --timeout             - Optional. Number of seconds to wait. Default 300 (5m)."
  exit 1
}

TIMEOUT=300
INTERVAL=2
TIMER_START=$SECONDS
WAIT_FOR_SUCCESSFUL_REQUESTS=3

while (( "$#" )); do
  case "$1" in
    -u|--url)
      SERVICE_URL=$2
      shift 2
      ;;
    -c|--cookie-jar)
      COOKIE_JAR=$2
      shift 2
      ;;
    -s|--successful-requests)
      WAIT_FOR_SUCCESSFUL_REQUESTS=$2
      shift 2
      ;;
    -t|--timeout)
      TIMEOUT=$2
      shift 2
      ;;
    -*|--*=)
      echo "Error: Unsupported option $1" >&2
      exit 1
      ;;
  esac
done

if [[ ! ${SERVICE_URL} ]]; then
 usage
fi

if [[ "${INTERVAL}" -gt "${TIMEOUT}" ]]; then
  INTERVAL=$TIMEOUT
fi

#
# Wait for endpoint to return a 200OK
#

SERVICE_CURL_RESULT=""
echo "> Waiting for ${SERVICE_URL} to return 200 OK (retrying every ${INTERVAL}s for ${TIMEOUT}s)"
limit=$(( ${TIMEOUT} / ${INTERVAL} ))
count=0
successful_requests=0
while : ; do
  printf "."
  if [ ! -z "${COOKIE_JAR}" ]; then
    SERVICE_CURL_RESULT=$(curl --cookie "${COOKIE_JAR}" --cookie-jar "${COOKIE_JAR}" -H 'Cache-Control: no-cache' -L -s -o /dev/null -w '%{http_code}' ${SERVICE_URL} || true)
  else
    SERVICE_CURL_RESULT=$(curl -H 'Cache-Control: no-cache' -L -s -o /dev/null -w '%{http_code}' ${SERVICE_URL} || true)
  fi

  if [[ "${SERVICE_CURL_RESULT}" -eq 200 ]]; then
    successful_requests=$[$successful_requests+1]
  elif [[ "${successful_requests}" -gt 0 ]]; then
    printf "\n[!] Warning: Previous request was successful, current request returned: ${SERVICE_CURL_RESULT}\n" >&2
    successful_requests=$[$successful_requests-1]
  fi

  if [[ "${successful_requests}" -ge "${WAIT_FOR_SUCCESSFUL_REQUESTS}" ]]; then
    printf "\n"
    break
  fi

  if [[ "${count}" -ge "${limit}" ]]; then
    printf "\n[!] Timeout waiting for Service to return 200 OK\n" >&2
    exit 1
  fi

  sleep ${INTERVAL}
  count=$[$count+1]
done

TIMER_DURATION=$(( SECONDS - TIMER_START ))

echo "> ${SERVICE_URL} returned ${successful_requests} successful requests in ${TIMER_DURATION}s"
