#!/usr/bin/env bash

dst=src/github.com/jvshahid
mkdir -p $dst
ln -f -s $PWD $dst/

. ./exports.sh

go get -v "launchpad.net/gocheck"

./run.sh test
