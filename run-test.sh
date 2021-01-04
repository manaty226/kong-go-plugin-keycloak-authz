#!/bin/bash

export CGO_CPPFLAGS="-Wno-error -Wno-nullability-completeness -Wno-expansion-to-defined -Wno-builtin-requires-header"

function functional_test() {
    cd ./test/functionalTest
    go test
}

function unit_test() {
    dirs=$(find $1 -maxdepth 1 -type d)

    for dir in $dirs;
    do
        if [[ $dir =~ .*/test || $dir =~ .*/.git ]]; then
            continue
        else
            cd $dir
            go test
            cd ../
        fi
    done

    cd $1
}


while [ "$#" != 0 ]
do
    if [ "$1" = "--unit" ]; then
        test_type=unit
    elif [ "$1" = "--functional" ]; then 
        test_type=functional
    fi
    shift
done

# sh ./test/start-keycloak.sh

cur_dir=$(pwd)

if [ "$test_type" = unit ]; then
    unit_test $cur_dir
else
    functional_test
fi

# sh ./test/stop-keycloak.sh > /dev/null




