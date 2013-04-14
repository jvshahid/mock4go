#!/usr/bin/env bash

dst=src/github.com/jvshahid
mkdir -p $dst
ln -f -s $PWD $dst/

bindir=$(dirname $(readlink -f $0))

. $bindir/exports.sh

go get -v "launchpad.net/gocheck"

rm -rf /tmp/gomock

export GOMOCK_TEST_ENV=gomock   # make sure we pass the environment properly

function test_package {
    rm -rf /tmp/gomock
    $bindir/run.sh $1
}

test_package test && test_package testc && test_package testnomock
