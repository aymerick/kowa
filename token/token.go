package token

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// Token represents a JWT token
type Token struct {
	Kind       string
	Value      interface{}
	Expiration int64
}

var signingKey []byte

// SigningKey returns the secret key setting
func SigningKey() []byte {
	if len(signingKey) == 0 {
		signingKey = []byte(viper.GetString("secret_key"))

		if len(signingKey) == 0 {
			panic("Signing key is missing")
		}
	}

	return signingKey
}

// NewToken instanciates a new Token
func NewToken(kind string, value interface{}) *Token {
	return &Token{
		Kind:  kind,
		Value: value,
	}
}

// SetExpirationTime sets expiration time on token
func (token *Token) SetExpirationTime(exp time.Time) {
	token.Expiration = exp.Unix()
}

// ExpirationTime returns tokenb expiration time
func (token *Token) ExpirationTime() time.Time {
	return time.Unix(int64(token.Expiration), 0)
}

// Expired returns true if token expired
func (token *Token) Expired() bool {
	return (token.Expiration != 0) && time.Now().After(token.ExpirationTime())
}

// Encode encodes token
func (token *Token) Encode() string {
	// create token
	t := jwt.New(jwt.SigningMethodHS256)

	t.Claims["k"] = token.Kind
	t.Claims["v"] = token.Value

	if token.Expiration != 0 {
		t.Claims["e"] = token.Expiration
	}

	// sign token
	result, err := t.SignedString(SigningKey())
	if err != nil {
		panic(err)
	}

	return result
}

// Decode decodes token
func Decode(encoded string) *Token {
	t, err := jwt.Parse(encoded, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return SigningKey(), nil
	})

	if err != nil {
		log.Printf("Failed to parse token: %v", err)
		return nil
	}

	// validate token
	if !t.Valid {
		log.Printf("Invalid token: %v", err)
		return nil
	}

	kind, isString := t.Claims["k"].(string)
	if !isString || (kind == "") {
		log.Printf("Invalid token kind: %v", err)
		return nil
	}

	result := NewToken(kind, t.Claims["v"])

	// get token expiry
	expiry, isFloat := t.Claims["e"].(float64)
	if isFloat {
		result.Expiration = int64(expiry)
	}

	return result
}
