#!/bin/bash

CONTAINER_TOOL=$1
COVERAGE_FILE=$2
TESTING_LIST=$(go list ./... | grep -v mocks | grep -v gen )

$CONTAINER_TOOL compose -f ./test/compose/dgraph.yml up --force-recreate --remove-orphans --detach

echo waiting for $CONTAINER_TOOL...
sleep 15
echo done waiting for $CONTAINER_TOOL, start testing...

go test -v -cover -coverprofile=$COVERAGE_FILE $TESTING_LIST

echo cleanup $CONTAINER_TOOL...
$CONTAINER_TOOL compose -f ./test/compose/dgraph.yml down