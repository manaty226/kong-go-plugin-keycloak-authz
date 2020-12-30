docker exec -i keycloak \
  /bin/bash -c "export JDBC_PARAMS=?currentSchema=keycloak_service && 
  /opt/jboss/keycloak/bin/standalone.sh \
  -Djboss.socket.binding.port-offset=100 \
  -Dkeycloak.migration.action=export \
  -Dkeycloak.migration.provider=singleFile \
  -Dkeycloak.migration.usersExportStrategy=REALM_FILE \
  -Dkeycloak.migration.file=/tmp/test-realm.json"

  docker cp keycloak:/tmp/test-realm.json ./test-realm.json