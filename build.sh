#!/usr/bin/env bash

. ./exports.sh

go list ./src/... | xargs go get -v
go build $args -v ./src/main/
