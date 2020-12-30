package keycloak

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"local.packages/token"
)

// IKeycloak is an interface of Keycloak
type IKeycloak interface {
	Protect(roles []string) bool
	Enforce(permissions []string) bool
}

// Keycloak is a basic keycloak structure
type Keycloak struct {
	Token        *token.Token
	ClientID     string
	Secret       string
	ServerURI    string
	Realm        string
	ResponseMode string
}

type permission struct {
	Rsid   string   `json:"rsid"`
	Scopes []string `json:"scopes"`
}

// Protect checks authentication and role of a received token
func (kc *Keycloak) Protect(roles []string) (hasPermit bool) {
	return true
}

// Enforce checks permissions of a received token
func (kc *Keycloak) Enforce(permissions []string) (hasPermit bool) {

	if len(permissions) == 0 {
		return true
	}

	permissionList := handlePermissions(permissions)

	return kc.checkPermissions(permissionList)
}

func handlePermissions(permissions []string) (permissionList map[string][]string) {
	var res map[string][]string = map[string][]string{}

	for _, p := range permissions {
		expected := strings.Split(p, ":")
		res[expected[0]] = []string{}
		if len(expected) == 2 {
			res[expected[0]] = append(res[expected[0]], expected[1])
		}
	}

	return res
}

func (kc *Keycloak) checkPermissions(expectedPermissions map[string][]string) (isAuthorized bool) {

	_, err := kc.getAuthorization(expectedPermissions)
	if err != nil {
		return false
	}

	// fmt.Printf("%v \n", permissions)

	return true
}

func (kc Keycloak) getAuthorization(expectedPermissions map[string][]string) (permissions []permission, err error) {
	values := url.Values{}
	values.Add("grant_type", "urn:ietf:params:oauth:grant-type:uma-ticket")
	values.Add("audience", kc.ClientID)
	values.Add("response_mode", kc.ResponseMode)
	createReqPermissions(&values, expectedPermissions)
	introspectionURL := kc.ServerURI + "/realms/" + kc.Realm + "/protocol/openid-connect/token"
	req, err := http.NewRequest(
		"Post",
		introspectionURL,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+kc.Token.GetJwt())

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error")
	}
	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("authorization error: %v", resp.StatusCode)
	}

	permissions = []permission{}
	if err := json.Unmarshal(resBody, &permissions); err != nil {
		return nil, fmt.Errorf("cannot read response json")
	}
	return permissions, nil
}

func createReqPermissions(values *url.Values, permissionList map[string][]string) {
	for resource := range permissionList {
		for _, scope := range permissionList[resource] {
			permission := resource + "#" + scope
			values.Add("permission", permission)
		}
	}
}

func hasPermission(permissions []permission, expectedPermission map[string][]string) (isIncluded bool) {
	for _, permission := range permissions {
		for _, scope := range permission.Scopes {
			if contains(expectedPermission[permission.Rsid], scope) {
				return true
			}
		}
	}
	return false
}

func contains(array []string, target string) (hasTarget bool) {
	for _, element := range array {
		if element == target {
			return true
		}
	}
	return false
}
