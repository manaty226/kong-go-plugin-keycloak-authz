#!/bin/bash

docker run -d \
--name keycloak \
-e KEYCLOAK_USER=admin \
-e KEYCLOAK_PASSWORD=admin \
-p 8080:8080 \
-v `pwd`/test/realm/test-realm.json:/config/test-realm.json \
jboss/keycloak \
-Djboss.bind.address.private=127.0.0.1 -Djboss.bind.address=0.0.0.0 \
-Dkeycloak.profile.feature.account_api=disabled \
-Dkeycloak.profile.feature.account2=disabled \
-Dkeycloak.migration.action=import \
-Dkeycloak.profile.feature.upload_scripts=enabled \
-Dkeycloak.migration.provider=singleFile \
-Dkeycloak.migration.file=/config/test-realm.json \
-Djboss.as.management.blocking.timeout=1200 \
-Dkeycloak.migration.strategy=OVERWRITE_EXISTING

counter=0
printf 'Starting keycloak'
until $(curl --output /dev/null --silent --head --fail http://localhost:8080/auth/realms/master/.well-known/openid-configuration); do
    printf '.'
    sleep 5
    if [[ "$counter" -gt 5 ]];then
        printf "failed to start keycloak server\n"
        docker rm keycloak
        exit 1
    fi
    counter=$((counter+1))
done

printf "\n"