package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/aymerick/kowa/commands"

	"gopkg.in/mgo.v2/bson"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SiteTestSuite struct {
	suite.Suite
}

// called before all tests
func (suite *SiteTestSuite) SetupSuite() {
	// Setup conf
	commands.ResetConf()
	commands.InitConf()
	// viper.Debug()

	// Change database
	SetDBName(fmt.Sprintf("%s_test", DBName()))
}

// called before each test
func (suite *SiteTestSuite) SetupTest() {
	// Reset database
	DB().DropDatabase()
}

// called after each test
func (suite *SiteTestSuite) TearDownTest() {
	// NOOP
}

// called after all tests
func (suite *SiteTestSuite) TearDownSuite() {
	// NOOP
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSiteTestSuite(t *testing.T) {
	suite.Run(t, new(SiteTestSuite))
}

//
// Tests
//

func (suite *SiteTestSuite) TestSite() {
	var err error

	t := suite.T()

	// Insert user
	user := User{
		Id:        bson.NewObjectId(),
		FirstName: "Jean-Claude",
		LastName:  "Trucmush",
		CreatedAt: time.Now(),
	}
	err = UsersCol().Insert(&user)
	assert.Nil(t, err)

	// Insert site
	site := Site{
		Id:          bson.NewObjectId(),
		UserId:      user.Id,
		CreatedAt:   time.Now(),
		Name:        "My site",
		Tagline:     "So powerfull !",
		Description: "You will be astonished by what my site is about",
	}
	err = SitesCol().Insert(&site)
	assert.Nil(t, err)

	// Count sites
	var c int
	c, err = SitesCol().Count()
	assert.Nil(t, err)

	assert.Equal(t, c, 1)

	// Fetch site
	var fetchedSite Site
	err = SitesCol().FindId(site.Id).One(&fetchedSite)
	assert.Nil(t, err)

	assert.Equal(t, fetchedSite.UserId, user.Id)
	assert.Equal(t, fetchedSite.Name, "My site")
	assert.Equal(t, fetchedSite.Tagline, "So powerfull !")
	assert.Equal(t, fetchedSite.Description, "You will be astonished by what my site is about")

	assert.NotNil(t, fetchedSite.Id)
	assert.NotNil(t, fetchedSite.CreatedAt)
}
