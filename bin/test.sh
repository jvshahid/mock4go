#!/usr/bin/env bash

dst=src/github.com/jvshahid
mkdir -p $dst
ln -f -s $PWD $dst/

cd `dirname $0`/..

. bin/exports.sh

go get -v "launchpad.net/gocheck"

rm -rf /tmp/gomock

echo "gopath: $GOPATH"

export GOMOCK_TEST_ENV=gomock   # make sure we pass the environment properly

function test_package {
    rm -rf /tmp/gomock
    go run gomock/gomock.go "$@"
}

if ! (test_package test && test_package testc && test_package testnomock); then
    echo "************************* TEST FAILED *******************************"
    exit 1
fi

exit 0
