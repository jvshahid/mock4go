#!/usr/bin/env bash

bindir=$(dirname $(readlink -f $0))

. $bindir/exports.sh

go run gomock/gomock.go "$@"
