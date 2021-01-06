package token

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
)

// Verify is verifying keycloak jwt signature
func verify(signed string, signature []byte, key *rsa.PublicKey) (isValid bool) {
	hasher := sha256.New()
	hasher.Write([]byte(signed))

	err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hasher.Sum(nil), signature)
	if err != nil {
		return false
	}

	return true
}

func getJWK() (key string) {
	return ""
}
