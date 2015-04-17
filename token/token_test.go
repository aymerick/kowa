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

	assert.Empty(t, token.Expiry)
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

	assert.Empty(t, token.Expiry)
	assert.Equal(t, "foo", token.Kind)
	assert.Equal(t, "bar", token.Value)
}

func (suite *TokenTestSuite) TestTokenExpiry() {
	t := suite.T()

	token := NewToken("foo", "bar")
	token.SetExpiration(time.Now().Add(time.Hour * 72))

	// encode
	encoded := token.Encode()
	assert.NotNil(t, encoded)

	// decode
	decoded := Decode(encoded)
	assert.NotNil(t, decoded)

	assert.Equal(t, token.Expiration(), decoded.Expiration())
	assert.Equal(t, "foo", decoded.Kind)
	assert.Equal(t, "bar", decoded.Value)
}
