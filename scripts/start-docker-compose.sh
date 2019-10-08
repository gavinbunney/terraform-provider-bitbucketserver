#!/bin/bash
set -e

echo "--> Downloading docker-compose"
curl -L https://github.com/docker/compose/releases/download/1.20.1/docker-compose-`uname -s`-`uname -m` > docker-compose
chmod +x docker-compose

echo "--> Starting docker-compose"
docker-compose up -d

echo "--> Wait for bitbucket to be ready"
bash -c 'while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' localhost:7990/status)" != "200" ]]; do sleep 5; done'
