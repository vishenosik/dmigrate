#!/bin/bash

COVERAGE_FILE=$1
TESTING_LIST=$(go list ./... | grep -v mocks | grep -v gen )

docker compose -f ./test/compose/docker-compose.yml up --force-recreate --remove-orphans --detach

echo waiting for docker...
sleep 10
echo done waiting for docker, start testing...

go test -v -cover -coverprofile=$COVERAGE_FILE $TESTING_LIST

echo cleanup docker...
docker compose -f ./test/compose/docker-compose.yml down