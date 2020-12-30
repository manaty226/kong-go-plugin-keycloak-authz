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
	Exp int64  `json:"exp"`
	Sub string `json:"sub"`
	Iss string `json:"iss"`
	Aud string `json:"aud"`
}

// NewToken is as constructor of Token
func NewToken(auth string) (token *Token, err error) {
	token = new(Token)
	if err := parseAccessToken(auth, token); err != nil {
		return nil, err
	}
	return token, nil
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
	token.Signature = splitAccessToken[2]

	return nil
}
