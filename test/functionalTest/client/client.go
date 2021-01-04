package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ClientConfig struct {
	ClientID     string
	ClientSecret string
	UserName     string
	Password     string
	ServerURL    string
}

type keycloakTokenRes struct {
	AccessToken string `json:"access_token"`
}

// Send is general method of client
func Send(serverURI string, method string, bodies map[string]string, headers map[string]string) (res *http.Response) {
	values := url.Values{}
	for key := range bodies {
		values.Add(key, bodies[key])
	}

	req, err := http.NewRequest(
		method,
		serverURI,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil
	}

	for key := range headers {
		req.Header.Set(key, headers[key])
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	fmt.Printf("%v \n", req)

	return resp
}

func GetUserToken(conf ClientConfig) (token string, err error) {
	values := url.Values{}

	values.Add("grant_type", "password")
	values.Add("client_id", conf.ClientID)
	values.Add("client_secret", conf.ClientSecret)
	values.Add("username", conf.UserName)
	values.Add("password", conf.Password)

	req, err := http.NewRequest(
		"Post",
		conf.ServerURL+"/realms/test/protocol/openid-connect/token",
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
