package token

import (
	"net/url"
	"time"

	"github.com/aymerick/kowa/models"
	"github.com/spf13/viper"
)

const (
	tokenAccountValidation = "account_validation"
)

// AccountActivationURL generate a token for user account activation
func AccountActivationURL(user *models.User) string {
	token := NewToken(tokenAccountValidation, user.ID)

	// token expires in 3 days
	token.SetExpirationTime(time.Now().Add(time.Hour * 72))

	// create URL
	endpoint, err := url.Parse(viper.GetString("service_url"))
	if err != nil {
		panic("Failed to parse service_url setting")
	}

	endpoint.Path += "/signup/validate"

	query := endpoint.Query()
	query.Set("token", token.Encode())
	endpoint.RawQuery = query.Encode()

	return endpoint.String()
}

// AccountValidationUser returns user id from token
func (token *Token) AccountValidationUser() string {
	if token.Kind != tokenAccountValidation {
		return ""
	}

	userID, ok := token.Value.(string)
	if !ok {
		return ""
	}

	return userID
}
