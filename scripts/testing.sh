#!/bin/bash

COVERAGE_FILE=$1
TESTING_LIST=$(go list ./... | grep -v mocks | grep -v gen )

go test -v -cover -coverprofile=$COVERAGE_FILE $TESTING_LIST