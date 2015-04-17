package token

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TokenTestSuite struct {
	suite.Suite
}

// called before all tests
func (suite *TokenTestSuite) SetupSuite() {
	viper.Set("secret_key", "my_so_secure_key")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTokenTestSuite(t *testing.T) {
	suite.Run(t, new(TokenTestSuite))
}

//
// Tests
//

func (suite *TokenTestSuite) TestTokenInstanciation() {
	t := suite.T()

	token := NewToken("foo", "bar")
	assert.NotNil(t, token)

	assert.Empty(t, token.Expiration)
	assert.False(t, token.Expired())

	assert.Equal(t, "foo", token.Kind)
	assert.Equal(t, "bar", token.Value)
}

func (suite *TokenTestSuite) TestTokenEncoding() {
	t := suite.T()

	encoded := NewToken("foo", "bar").Encode()
	assert.NotNil(t, encoded)
}

func (suite *TokenTestSuite) TestTokenDecoding() {
	t := suite.T()

	encoded := NewToken("foo", "bar").Encode()

	token := Decode(encoded)
	assert.NotNil(t, token)

	assert.Empty(t, token.Expiration)
	assert.False(t, token.Expired())

	assert.Equal(t, "foo", token.Kind)
	assert.Equal(t, "bar", token.Value)
}

func (suite *TokenTestSuite) TestTokenExpiration() {
	t := suite.T()

	token := NewToken("foo", "bar")
	token.SetExpirationTime(time.Now().Add(time.Hour * 72))

	// encode
	encoded := token.Encode()
	assert.NotNil(t, encoded)

	// decode
	decoded := Decode(encoded)
	assert.NotNil(t, decoded)

	assert.Equal(t, token.ExpirationTime(), decoded.ExpirationTime())
	assert.False(t, token.Expired())

	assert.Equal(t, "foo", decoded.Kind)
	assert.Equal(t, "bar", decoded.Value)
}

func (suite *TokenTestSuite) TestTokenExpired() {
	t := suite.T()

	token := NewToken("foo", "bar")
	token.SetExpirationTime(time.Now().Add(time.Hour * -1))

	assert.NotEmpty(t, token.Expiration)
	assert.True(t, token.Expired())
}
