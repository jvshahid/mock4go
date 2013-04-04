#!/usr/bin/env bash

orig=`readlink -f $1`

cd `dirname $0`

if [ $# -ne 1 ]; then
    echo "Usage: $0 <file>"
    exit 1
fi

. bin/exports.sh

git_root_dir=$(git rev-parse --show-toplevel)
orig=${orig/$git_root_dir\//}
dirname=`dirname $orig`
filename=`basename $orig`
package=${dirname/*src\//}
echo $orig

function is_test {
    if [[ "$orig" == *_test.go ]]; then
        return 0
    else
        return 1
    fi
}

function flymake_regular {
    # go build $orig
    go build $orig 2>&1 | grep $orig | sed "s/.*$filename/$filename/g"
}

function flymake_test {
    go test -c $package 2>&1 | grep $orig | sed "s/.*$filename/$filename/g"
}

is_test
if [ $? -eq 0 ]; then
    output=$(flymake_test)
else
    output=$(flymake_regular)
fi
if [ "x$output" == "x" ]; then
    exit 0
else
    echo "$output"
    exit 1
fi
