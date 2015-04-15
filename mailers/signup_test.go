package mailers

import (
	"bytes"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"path"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aymerick/kowa/core"
	"github.com/aymerick/kowa/helpers"
	"github.com/aymerick/kowa/models"
)

type SignupTestSuite struct {
	suite.Suite
}

// called before all tests
func (suite *SignupTestSuite) SetupSuite() {
	core.LoadLocales()

	SetTemplatesDir(path.Join(helpers.WorkingDir(), "templates"))

	viper.Set("smtp_from", "test@test.com")
	viper.Set("service_name", "My Service")
	// viper.Set("service_logo", "http://www.myservice.bar/logo.png")
	viper.Set("service_url", "http://www.myservice.bar")
	viper.Set("service_postal_address", "2 quality street - 1337 GGCity - RoxLand")
	viper.Set("service_copyright_notice", "Copyright @ 2015 AceOfBase - All rights reserved")
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

	// Instanciate user
	user := &models.User{
		Id:        "trucmush",
		Email:     "trucmush@wanadoo.fr",
		FirstName: "Jean-Claude",
		LastName:  "Trucmush",
		CreatedAt: time.Now(),
		Lang:      "en",
	}

	sender := NewSender(NewSignupMailer(user))
	sender.SetNoop(true)
	sender.SetSMTPConf(&SMTPConf{
		From: "Marcel Belivo <surprise@surpri.se>",
		Host: "pantoute.com",
		Port: 561,
		User: "jeanmich",
		Pass: "troudku",
	})

	// log.Printf("HTML MAIL:\n\n%v", sender.Content(TPL_HTML))

	email := sender.NewEmail()
	assert.NotNil(t, email)

	// log.Printf("HTML INLINE MAIL:\n\n%v", string(email.HTML))

	errSend := sender.Send()
	assert.Nil(t, errSend)

	// check mail generation
	rawMail, errGen := email.Bytes()
	assert.Nil(t, errGen)

	// parse generated mail
	msg, errRead := mail.ReadMessage(bytes.NewBuffer(rawMail))
	assert.Nil(t, errRead)

	// check headers
	expectedHeaders := map[string]string{
		"To":      "Jean-Claude Trucmush <trucmush@wanadoo.fr>",
		"From":    "Marcel Belivo <surprise@surpri.se>",
		"Subject": "Activate your My Service account.",
	}

	for header, expected := range expectedHeaders {
		val := msg.Header.Get(header)
		assert.Equal(t, expected, val)
	}

	ct := msg.Header.Get("Content-type")
	mt, params, errMt := mime.ParseMediaType(ct)
	assert.Nil(t, errMt)
	assert.Equal(t, "multipart/mixed", mt)

	mixedBoundary := params["boundary"]
	assert.NotEmpty(t, mixedBoundary)

	mixed := multipart.NewReader(msg.Body, mixedBoundary)
	textPart, errPart := mixed.NextPart()
	assert.Nil(t, errPart)

	mt, params, errMt = mime.ParseMediaType(textPart.Header.Get("Content-type"))
	assert.Nil(t, errMt)
	assert.Equal(t, "multipart/alternative", mt)

	mpReader := multipart.NewReader(textPart, params["boundary"])
	partText, errN := mpReader.NextPart()
	assert.Nil(t, errN)

	plainText, errPT := ioutil.ReadAll(partText)
	assert.Nil(t, errPT)

	textStr := string(plainText)
	// log.Printf("plainText:\n\n%s", textStr)
	assert.Regexp(t, `Thanks for joining My Service.`, textStr)
	assert.Regexp(t, `Just one more step...`, textStr)
	assert.Regexp(t, `Activate Your Account`, textStr)

	partHtml, errH := mpReader.NextPart()
	assert.Nil(t, errH)

	htmlContent, errHC := ioutil.ReadAll(partHtml)
	assert.Nil(t, errHC)

	htmlStr := string(htmlContent)
	// log.Printf("htmlContent:\n\n%s", htmlStr)
	assert.Regexp(t, `<title>Activate your My Service account\.</title>`, htmlStr)
	assert.Regexp(t, `Click the button below to activate your account\.`, htmlStr)
	assert.Regexp(t, `Activate account`, htmlStr)
}
