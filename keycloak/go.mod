module keycloak

go 1.15

replace local.packages/token => ../token

replace local.packages/keycloak => ./

require (
	local.packages/keycloak v0.0.0-00010101000000-000000000000
	local.packages/token v0.0.0-00010101000000-000000000000
)
