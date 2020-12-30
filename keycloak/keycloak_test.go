package keycloak

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/manaty226/kong-go-plugin-keycloak-authz/token"
)

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

func getToken() (token string, err error) {
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

func Test_Keycloak(t *testing.T) {
	accessToken, err := getToken()
	if err != nil {
		t.Errorf("cannot get token.")
	}

	newToken, _ := token.NewToken("Bearer " + accessToken)
	kc := Keycloak{
		Token:        newToken,
		ClientID:     ClientID,
		Secret:       ClientSecret,
		ServerURI:    ServerURL,
		Realm:        Realm,
		ResponseMode: ResponseMode,
	}

	if !kc.Enforce([]string{"testResource:testScope"}) {
		t.Errorf("not authenticated")
	}

}
