#!/bin/bash

export CGO_CPPFLAGS="-Wno-error -Wno-nullability-completeness -Wno-expansion-to-defined -Wno-builtin-requires-header"

root_dir=$(pwd)
sh ./test/start-keycloak.sh

dirs=$(find ./* -maxdepth 0 -type d)

for dir in $dirs;
do
    if [ $dir != "./test" ]; then
        echo $dir
        cd $dir
        go test
        cd ../
    fi
done

cd $root_dir

sh ./test/stop-keycloak.sh > /dev/null