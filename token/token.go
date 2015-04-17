package token

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type Token struct {
	Kind   string
	Value  interface{}
	Expiry int64
}

var signingKey []byte

func SigningKey() []byte {
	if len(signingKey) == 0 {
		signingKey = []byte(viper.GetString("secret_key"))

		if len(signingKey) == 0 {
			panic("Signing key is missing")
		}
	}

	return signingKey
}

func NewToken(kind string, value interface{}) *Token {
	return &Token{
		Kind:  kind,
		Value: value,
	}
}

func (token *Token) SetExpiration(exp time.Time) {
	token.Expiry = exp.Unix()
}

func (token *Token) Expiration() time.Time {
	return time.Unix(int64(token.Expiry), 0)
}

func (token *Token) Encode() string {
	// create token
	t := jwt.New(jwt.SigningMethodHS256)

	t.Claims["k"] = token.Kind
	t.Claims["v"] = token.Value

	if token.Expiry != 0 {
		t.Claims["e"] = token.Expiry
	}

	// sign token
	result, err := t.SignedString(SigningKey())
	if err != nil {
		panic(err)
	}

	return result
}

func Decode(encoded string) *Token {
	t, err := jwt.Parse(encoded, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("Unexpected signing method: %v", t.Header["alg"]))
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
		result.Expiry = int64(expiry)
	}

	return result
}
