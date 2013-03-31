#!/usr/bin/env bash

. ./exports.sh

go get -v "launchpad.net/gocheck"

./run.sh test
