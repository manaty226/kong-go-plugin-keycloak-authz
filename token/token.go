package token

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// IToken is an Interface type of Token
type IToken interface {
	GetJwt() string
	IsExpired() bool
}

// Token is a class which contains parsed access token
type Token struct {
	jwt       string
	Header    jwtHeader
	Content   jwtContent
	Signature string
}

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type jwtContent struct {
	Exp            int64            `json:"exp"`
	Sub            string           `json:"sub"`
	Iss            string           `json:"iss"`
	Aud            string           `json:"aud"`
	RealmAccess    roles            `json:"realm_access"`
	ResourceAccess map[string]roles `json:"resource_access"`
}

type roles struct {
	Roles []string `json:"roles"`
}

// NewToken is as constructor of Token
func NewToken(auth string) (token *Token, err error) {
	token = new(Token)
	if err := parseAccessToken(auth, token); err != nil {
		return nil, err
	}
	return token, nil
}

// HasRole checks the token contains arg roles
func (t *Token) HasRole(roles []string, clientID string) (hasRole bool) {
	var splitRole []string

	for _, role := range roles {
		splitRole = strings.Split(role, ":")
		if len(splitRole) == 1 {
			return t.hasApplicationRole(clientID, splitRole[0])
		} else if splitRole[0] == "realm" {
			return t.hasRealmRole(splitRole[1])
		}
	}
	return t.hasApplicationRole(splitRole[0], splitRole[1])
}

func (t *Token) hasRealmRole(role string) (hasRole bool) {

	for _, roleName := range RealmAccess {
		if role == roleName {
			return true
		}
	}

	return false
}

func (t *Token) hasApplicationRole(app string, role string) (hasRole bool) {

	if appRoles, hasKey := t.Content.ResourceAccess[app]; hasKey {
		for _, roleName := range appRoles {
			if role == roleName {
				return true
			}
		}
	}
	return false
}

// GetJwt is to get raw jwt token
func (t *Token) GetJwt() (jwt string) {
	return t.jwt
}

// IsExpired checks if the token is expired
func (t *Token) IsExpired() (isExpired bool) {
	return t.Content.Exp < time.Now().Unix()
}

func parseAccessToken(auth string, token *Token) (err error) {
	splitAuthHeader := strings.Split(auth, "Bearer ")
	if splitAuthHeader[1] == "" {
		return fmt.Errorf("token not contained")
	}
	token.jwt = splitAuthHeader[1]
	accessToken := strings.TrimSpace(splitAuthHeader[1])
	splitAccessToken := strings.Split(accessToken, ".")
	if len(splitAccessToken) != 3 {
		return fmt.Errorf("invalid jwt format")
	}

	header, err := base64.RawURLEncoding.DecodeString(splitAccessToken[0])
	if err := json.Unmarshal(header, &token.Header); err != nil {
		return err
	}
	content, err := base64.RawURLEncoding.DecodeString(splitAccessToken[1])
	if err := json.Unmarshal(content, &token.Content); err != nil {
		return err
	}

	fmt.Printf("%v \n", string(content))

	fmt.Printf("token is %v \n", token.Content.ResourceAccess)

	token.Signature = splitAccessToken[2]

	return nil
}
