#!/bin/bash

export CGO_CPPFLAGS="-Wno-error -Wno-nullability-completeness -Wno-expansion-to-defined -Wno-builtin-requires-header"

cur_dir=$(pwd)
sh ./test/start-keycloak.sh

dirs=$(find ./lib/* -type d)

for dir in $dirs;
do
    echo $dir
    cd $dir
    go test
    cd ../../
done

cd $cur_dir

sh ./test/stop-keycloak.sh > /dev/null