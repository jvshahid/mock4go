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
    go run gomock/gomock.go "$@"
    if [ -d $TMPDIR/gomock ];then
        echo "gomock didn't cleanup the gomock directory. WTF"
        exit 1
    fi
}

if [ "x$TMPDIR" == "x" ]; then
    TMPDIR=/tmp
fi
destination=$TMPDIR/gomock_$BASHPID

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
    echo "Cannot find $destination although we ran gomock with -i and -k"
    echo "************************* TEST FAILED *******************************"
    exit 1
fi

cleanup
exit 0
