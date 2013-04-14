#!/usr/bin/env bash

dst=src/github.com/jvshahid
mkdir -p $dst
ln -f -s $PWD $dst/

bindir=$(dirname $(readlink -f $0))

. $bindir/exports.sh

go get -v "launchpad.net/gocheck"

export GOMOCK_TEST_ENV=gomock   # make sure we pass the environment properly
$bindir/run.sh test && $bindir/run.sh testc
