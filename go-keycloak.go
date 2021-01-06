package main

import (
	"github.com/manaty226/kong-go-plugin-keycloak-authz/keycloak"
	"github.com/manaty226/kong-go-plugin-keycloak-authz/response"
	"github.com/manaty226/kong-go-plugin-keycloak-authz/token"

	"github.com/Kong/go-pdk"
)

// Config is a config structure of this plugin
type Config struct {
	Mode      string   `json:"mode"`
	ClientID  string   `json:"client_id"`
	Secret    string   `json:"secret"`
	ServerURI string   `json:"server_uri"`
	Realm     string   `json:"realm"`
	Rules     []string `json:"rules"`
}

// New is required for a kong-plugin
func New() interface{} {
	return &Config{}
}

// Access is main function of this plugin
func (conf Config) Access(kong *pdk.PDK) {

	auth, err := kong.Request.GetHeader("Authorization")
	if err != nil {
		kong.Log.Err("No access token.")
		statusCode, respBody, headers := response.AuthErrorResponse()
		kong.Response.SetStatus(statusCode)
		kong.Response.Exit(statusCode, respBody, headers)
		return
	}

	t, err := token.NewToken(auth)
	if err != nil {
		kong.Log.Err("parse error: ", err.Error())
		statusCode, respBody, headers := response.AuthErrorResponse()
		kong.Response.SetStatus(statusCode)
		kong.Response.Exit(statusCode, respBody, headers)
		return
	}

	kc := keycloak.Keycloak{
		Token:        t,
		ClientID:     conf.ClientID,
		Secret:       conf.Secret,
		ServerURI:    conf.ServerURI,
		Realm:        conf.Realm,
		ResponseMode: "permissions",
	}

	isAuthorized := false
	if conf.Mode == "Enforce" {
		isAuthorized = kc.Enforce(conf.Rules)
	} else if conf.Mode == "Protect" {
		isAuthorized = kc.Protect(conf.Rules)
	} else {
		kong.Log.Err("mode setting error. All accesses are denied.")
	}

	if !isAuthorized {
		statusCode, respBody, headers := response.AuthErrorResponse()
		kong.Response.SetStatus(statusCode)
		kong.Response.Exit(statusCode, respBody, headers)
	}

}
