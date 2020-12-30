module github.com/manaty226/kong-go-plugin-keycloak-authz 

go 1.13

require (
	github.com/Kong/go-pdk v0.5.0
	github.com/Kong/go-pluginserver v0.6.0 // indirect
)

replace local.packages/response => ./response
replace local.packages/token => ./token
replace local.packages/keycloak => ./keycloak