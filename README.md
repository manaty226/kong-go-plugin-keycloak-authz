# kong-go-plugin-keycloak-authz

## Description
This repository contains a kong plugin which works as policy enforcement point
to check access token provided by Keycloak.
Similar to [keycloak-nodejs-connect](https://github.com/keycloak/keycloak-nodejs-connect), role base access control (RBAC) and policy base access control (PBAC) are realized by this plugin.
## Quick start
Kong with this plugin and keycloak can be run by docker-compose as follows.
````
docker-compose -d up kong keycloak
````

## Configuration

| name | description |
| ---  | ---|
| mode | If you set as `Protect`, this plugin works as RABC. If `Enforce`, working as PBAC. |
| client_id | client id of keycloak realm. |
| secret | client secret of keycloak realm |
| server uri | keycloak server uri |
| realm | The name of realm |
| rules | When the mode is set as `Protect`, approved role names have to be set in this parameter. When the mode is set as `Enforce`, resource names and scopes have to be set. | 