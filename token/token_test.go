package token

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

var jwt string = "Bearer " +
	"eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJYTy1sY1VuYWFtOFg2X2FCXzhaVUx3YTh0cUs0X1NNUHhGaDRsRVg5TzR3In0.eyJleHAiOjE2MDg2MzYyOTEsImlhdCI6MTYwODYzNjIzMSwianRpIjoiNzE0NWFiMWYtZTJhNS00MTZkLWJlZTQtMmUyNjU1OWZjNDM1IiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL2F1dGgvcmVhbG1zL21hc3RlciIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiJkNWU5OTRkZi1jYTZjLTQxYTUtYjM0YS03Y2EwNDc1NzQwYWIiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJ0ZXN0LWNsaWVudCIsInNlc3Npb25fc3RhdGUiOiJkMzUyNDhiYi0xODRlLTQ2ZjQtOGY2OC01ZTQ1YTcwNjRjNjkiLCJhY3IiOiIxIiwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6Im9wZW5pZCBlbWFpbCBwcm9maWxlIiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJ0ZXN0In0.U2aih1sKy5fy0CcrPfXwvO94_UI-mKsXQ34rlNrKJgseAEVtn_fpdA2UO9JEjqZ6YDfuMB4DN-nBT6TYjSwrYlBVGl4ofihRY_4VjhzdtF726GvRyNRRRmslrSf6z6aycclwqms8qOi67C7Pn2QKhhbT8zckcKQQz87B2H3cwOhCfbCcGtYdRbICs7YForX6h7ahpvP79qzTk5-5omEgHl8J6NTs9ykPPU7okqpd9jP8RCjDCYTPsqYcTxFckRjSZeJr3J7hu0qGp1Z01fC7Ppgwlm2jTGSiPGp5LNWzdchNLKRJb77ogROM32wvoz1MaOhuMzk9Dx56QfwDBr1E0A"

type keycloakTokenRes struct {
	AccessToken string `json:"access_token"`
}

type tokenReqBody struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

var ClientID string = "test-client"
var ClientSecret string = "ba15024a-725b-45ea-bac2-c70332e4c4d7"
var ServerURL string = "http://localhost:8080/auth"
var Realm string = "test"
var ResponseMode string = "permissions"

func getUserToken() (token string, err error) {
	values := url.Values{}

	values.Add("grant_type", "password")
	values.Add("client_id", ClientID)
	values.Add("client_secret", ClientSecret)
	values.Add("username", "test")
	values.Add("password", "test")

	req, err := http.NewRequest(
		"Post",
		ServerURL+"/realms/test/protocol/openid-connect/token",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	resBody, _ := ioutil.ReadAll(resp.Body)

	kcRes := new(keycloakTokenRes)
	if err := json.Unmarshal(resBody, &kcRes); err != nil {
		return "", fmt.Errorf("cannot read response json")
	}

	return kcRes.AccessToken, nil
}

func Test_ParseToken(t *testing.T) {
	accessToken, err := NewToken(jwt)

	if err != nil {
		t.Errorf(err.Error())
	}

	expected := "account"

	if accessToken.Content.Aud != expected {
		t.Errorf("got: %v\n want: %v", accessToken.Content.Aud, expected)
	}
}

func Test_GetJwt(t *testing.T) {
	accessToken, err := NewToken(jwt)
	if err != nil {
		t.Errorf(err.Error())
	}

	gotJwt := "Bearer " + accessToken.GetJwt()

	if gotJwt != jwt {
		t.Errorf("cannot get correct jwt: got %v", gotJwt)
	}
}

func Test_IsExpired(t *testing.T) {
	accessToken, err := NewToken(jwt)
	if err != nil {
		t.Errorf((err.Error()))
	}

	if !accessToken.IsExpired() {
		t.Errorf("cannot check expired.")
	}
}

func Test_EmptyToken(t *testing.T) {
	_, err := NewToken("Bearer ")

	if err == nil {
		t.Errorf("empty token cannot be correctly treat")
	}
}

func Test_Sign(t *testing.T) {
	token, _ := getUserToken()
	accessToken, _ := NewToken("Bearer " + token)

	if accessToken.IsValidSignature() != true {
		t.Errorf("invalid Signature")
	}
}
