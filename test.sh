#!/bin/bash

# Run all unit tests.
echo "Running unit tests .."
go test github.com/yashmurty/go-web-crawler -cover -count=1

if [[ ! -z "$1" ]]; then
  # Run all e2e tests.
  echo "Running e2e tests .."
  go test github.com/yashmurty/go-web-crawler/e2e -cover -count=1 -timeout 120s
fi
echo "Done."
