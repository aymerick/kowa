package token

import (
	"strings"
	"testing"
	"time"

	"github.com/aymerick/kowa/models"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TokenUserTestSuite struct {
	suite.Suite
}

// called before all tests
func (suite *TokenUserTestSuite) SetupSuite() {
	viper.Set("secret_key", "my_so_secure_key")
	viper.Set("service_url", "http://www.myservice.bar")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTokenUserTestSuite(t *testing.T) {
	suite.Run(t, new(TokenUserTestSuite))
}

//
// Tests
//

func (suite *TokenUserTestSuite) TestAccountActivationURL() {
	t := suite.T()

	user := &models.User{
		ID:        "trucmush",
		Email:     "trucmush@wanadoo.fr",
		FirstName: "Jean-Claude",
		LastName:  "Trucmush",
		CreatedAt: time.Now(),
		Lang:      "en",
	}

	url := AccountActivationURL(user)
	assert.NotNil(t, url)

	expectedPrefix := "http://www.myservice.bar/signup/validate?token="

	assert.True(t, strings.HasPrefix(url, expectedPrefix))
	encoded := url[len(expectedPrefix):len(url)]

	decoded := Decode(encoded)
	assert.NotNil(t, decoded)

	assert.NotNil(t, decoded.ExpirationTime())
	assert.Equal(t, tokenAccountValidation, decoded.Kind)
	assert.Equal(t, user.ID, decoded.Value)
}
