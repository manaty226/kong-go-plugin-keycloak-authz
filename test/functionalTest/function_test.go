package main

import (
	"functionalTest/client"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func Test_Keycloak(t *testing.T) {

	token := getUserToken()

	res := client.Send(
		"http://kong:8081",
		"GET",
		map[string]string{},
		map[string]string{
			"Authorization": "Bearer " + token,
		},
	)

	if res.StatusCode != 200 {
		t.Errorf("expected: passed, actual: invalid access")
	}
}

func getUserToken() (token string) {

	conf := client.ClientConfig{
		ClientID:     "test-client",
		ClientSecret: "ba15024a-725b-45ea-bac2-c70332e4c4d7",
		UserName:     "test",
		Password:     "test",
		ServerURL:    "http://keycloak:8080/auth",
	}

	accessToken, _ := client.GetUserToken(conf)
	return accessToken
}
