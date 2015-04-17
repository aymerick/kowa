package token

import (
	"net/url"
	"time"

	"github.com/aymerick/kowa/models"
	"github.com/spf13/viper"
)

const (
	TOKEN_ACCOUNT_VALIDATION = "account_validation"
)

// Generate a token for user account activation
func AccountActivationUrl(user *models.User) string {
	token := NewToken(TOKEN_ACCOUNT_VALIDATION, user.Id)

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

func (token *Token) AccountValidationUser() string {
	if token.Kind != TOKEN_ACCOUNT_VALIDATION {
		return ""
	}

	userId, ok := token.Value.(string)
	if !ok {
		return ""
	}

	return userId
}
