package token

import (
	"crypto/rsa"
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
	Signature []byte
	signed    string
}

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
	Kid string `json:"kid"`
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
	for _, role := range roles {
		splitRole := strings.Split(role, ":")
		if len(splitRole) == 1 {
			if t.hasApplicationRole(clientID, splitRole[0]) {
				return true
			}
		} else if splitRole[0] == "realm" {
			if t.hasRealmRole(splitRole[1]) {
				return true
			}
		} else {
			if t.hasApplicationRole(splitRole[0], splitRole[1]) {
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

// IsValidSignature check signature of jwt
func (t *Token) IsValidSignature(retrieveKey func() *rsa.PublicKey) (isValid bool) {
	if retrieveKey() == nil {
		return false
	}
	return verify(t.signed, t.Signature, retrieveKey())
}

func (t *Token) hasRealmRole(role string) (hasRole bool) {
	for _, tokenRole := range t.Content.RealmAccess.Roles {
		if role == tokenRole {
			return true
		}
	}
	return false
}

func (t *Token) hasApplicationRole(app string, role string) (hasRole bool) {
	if resourceAccess, hasKey := t.Content.ResourceAccess[app]; hasKey {
		for _, tokenRole := range resourceAccess.Roles {
			if role == tokenRole {
				return true
			}
		}
	}
	return false
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

	signature, err := base64.RawURLEncoding.DecodeString(splitAccessToken[2])
	if err != nil {
		return err
	}
	token.Signature = signature
	token.signed = splitAccessToken[0] + "." + splitAccessToken[1]

	return nil
}
