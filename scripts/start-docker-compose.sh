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
${DIR}/docker-compose up -d --build

echo "--> Wait for bitbucket to be ready"
#If the BITBUCKET_SERVER environment variable is not set then use http://localhost:7990
bash ${DIR}/wait-for-url.sh --url ${BITBUCKET_SERVER-http://localhost:7990}/status --timeout 600
