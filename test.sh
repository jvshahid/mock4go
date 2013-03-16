#!/usr/bin/env bash

. ./exports.sh

go get -v "launchpad.net/gocheck"

go test test "$@"
