#!/bin/bash

export CGO_CPPFLAGS="-Wno-error -Wno-nullability-completeness -Wno-expansion-to-defined -Wno-builtin-requires-header"

function start_keycloak() {
    docker-compose up -d keycloak
    until $(curl --output /dev/null --silent --head --fail http://localhost:8080/auth/realms/master/.well-known/openid-configuration); do
        printf '.'
        sleep 5
        if [[ "$counter" -gt 5 ]]; then
            printf "failed to start keycloak server\n"
            docker rm keycloak
            exit 1
        fi
        counter=$((counter+1))
    done
}

function functional_test() {
    start_keycloak
    docker-compose up -d kong
    docker-compose up client
    docker-compose down
}

function unit_test() {
    start_keycloak

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

docker-compose build kong

cur_dir=$(pwd)
cd ./test/

if [ "$test_type" = unit ]; then
    unit_test $cur_dir
else
    functional_test
fi




