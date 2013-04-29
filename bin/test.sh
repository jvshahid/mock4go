#!/usr/bin/env bash

dst=src/github.com/jvshahid
mkdir -p $dst
ln -f -s $PWD $dst/

cd `dirname $0`/..

. bin/exports.sh

go get -v "launchpad.net/gocheck"

rm -rf /tmp/mock4go

echo "gopath: $GOPATH"

export MOCK4GO_TEST_ENV=mock4go   # make sure we pass the environment properly

function test_package {
    go run mock4go/mock4go.go "$@"
    status=$?
    if [ -d $TMPDIR/mock4go ];then
        echo "mock4go didn't cleanup the mock4go directory. WTF"
        exit 1
    fi
    return $status
}

if [ "x$TMPDIR" == "x" ]; then
    TMPDIR=/tmp
fi
destination=$TMPDIR/mock4go_$BASHPID

function cleanup {
    rm -rf $destination
}

trap cleanup EXIT

if ! (test_package test && test_package testc && test_package testnomock && \
        test_package -i -k -d $destination test_failing); then
    echo "************************* TEST FAILED *******************************"
    exit 1
fi

trap - EXIT
if [ ! -d $destination ];then
    echo "Cannot find $destination although we ran mock4go with -i and -k"
    echo "************************* TEST FAILED *******************************"
    exit 1
fi

cleanup
exit 0
