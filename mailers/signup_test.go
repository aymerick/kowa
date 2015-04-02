package mailers

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aymerick/kowa/core"
	"github.com/aymerick/kowa/models"
)

type SignupTestSuite struct {
	suite.Suite
}

// called before all tests
func (suite *SignupTestSuite) SetupSuite() {
	core.LoadLocales()

	viper.Set("smtp_from", "test@test.com")
}

// called before each test
func (suite *SignupTestSuite) SetupTest() {
	// NOOP
}

// called after each test
func (suite *SignupTestSuite) TearDownTest() {
	// NOOP
}

// called after all tests
func (suite *SignupTestSuite) TearDownSuite() {
	// NOOP
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSignupTestSuite(t *testing.T) {
	suite.Run(t, new(SignupTestSuite))
}

//
// Tests
//

func (suite *SignupTestSuite) TestSignup() {
	t := suite.T()

	// Insert user
	user := &models.User{
		Id:        "trucmush",
		Email:     "trucmush@wanadoo.fr",
		FirstName: "Jean-Claude",
		LastName:  "Trucmush",
		CreatedAt: time.Now(),
		Lang:      "en",
	}

	sender := NewSender(NewSignupMailer(user))

	mail := sender.NewEmail()
	assert.NotNil(t, mail)

	_, err := mail.Bytes()
	assert.Nil(t, err)

	// assert.Equal(t, "", string(rawMail))
}
