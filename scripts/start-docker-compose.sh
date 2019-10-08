#!/bin/bash
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

cd ${DIR}

if [ ! -f ${DIR}/docker-compose ]; then
  echo "--> Downloading docker-compose"
  curl -L https://github.com/docker/compose/releases/download/1.20.1/docker-compose-`uname -s`-`uname -m` > ${DIR}/docker-compose
  chmod +x ${DIR}/docker-compose
fi

echo "--> Starting docker-compose"
${DIR}/docker-compose up -d

echo "--> Wait for bitbucket to be ready"
bash -c 'while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' localhost:7990/status)" != "200" ]]; do sleep 5; done'
