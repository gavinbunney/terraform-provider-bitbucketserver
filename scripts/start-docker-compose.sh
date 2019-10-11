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
bash ${DIR}/wait-for-url.sh --url http://localhost:7990/status --timeout 600

echo "--> Install a plugin for tests to reference licensing etc"
if [ ! -f ${DIR}/docker-compose-test-plugin.jar ]; then
  curl -L -o ${DIR}/docker-compose-test-plugin.jar https://marketplace.atlassian.com/download/apps/1211185/version/400500101
fi

url="http://admin:admin@localhost:7990/rest/plugins/1.0/"; \
token=$(curl  -H 'X-Atlassian-Token: no-check' -sI "$url?os_authType=basic" | grep upm-token | cut -d: -f2- | tr -d '[[:space:]]'); \
curl -H 'X-Atlassian-Token: no-check' -XPOST "$url?token=$token" -F plugin=@${DIR}/docker-compose-test-plugin.jar
sleep 10

echo "--> Apply a timebomb license for the plugin"
PLUGIN_LICENSE="AAABCA0ODAoPeNpdj01PwkAURffzKyZxZ1IyUzARkllQ24gRaQMtGnaP8VEmtjPNfFT59yJVFyzfubkn796Ux0Bz6SmbUM5nbDzj97RISxozHpMUnbSq88poUaLztFEStUN6MJZ2TaiVpu/YY2M6tI6sQrtHmx8qd74EZ+TBIvyUU/AoYs7jiE0jzknWQxMuifA2IBlUbnQ7AulVjwN9AaU9atASs69O2dNFU4wXJLc1aOUGw9w34JwCTTZoe7RPqUgep2X0Vm0n0fNut4gSxl/Jcnj9nFb6Q5tP/Ueu3L+0PHW4ghZFmm2zZV5k6/95CbR7Y9bYGo/zGrV3Ir4jRbDyCA6vt34DO8p3SDAsAhQnJjLD5k9Fr3uaIzkXKf83o5vDdQIUe4XequNCC3D+9ht9ZYhNZFKmnhc=X02dh"
PLUGIN_LICENSE_ENDPOINT="http://admin:admin@localhost:7990/rest/plugins/1.0/nl.stefankohler.stash.stash-notification-plugin-key/license?os_authType=basic"

curl -H 'X-Atlassian-Token: no-check' -H 'Content-Type: application/vnd.atl.plugins+json' \
    -X PUT -d "{\"rawLicense\": \"${PLUGIN_LICENSE}\"}" ${PLUGIN_LICENSE_ENDPOINT}
