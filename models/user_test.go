package models

import (
	"encoding/json"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	db *DBSession
}

// called before all tests
func (suite *UserTestSuite) SetupSuite() {
	// setup db
	suite.db = NewTestDBSession()
	suite.db.SetDBName(TEST_DBNAME)
}

// called before each test
func (suite *UserTestSuite) SetupTest() {
	// Reset database
	suite.db.DB().DropDatabase()
}

// called after each test
func (suite *UserTestSuite) TearDownTest() {
	// NOOP
}

// called after all tests
func (suite *UserTestSuite) TearDownSuite() {
	// NOOP
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

//
// Tests
//

func (suite *UserTestSuite) TestUsers() {
	var err error

	t := suite.T()

	// Insert users
	err = suite.db.UsersCol().Insert(
		&User{
			ID:        "trucmush",
			Email:     "trucmush@wanadoo.fr",
			FirstName: "Jean-Claude",
			LastName:  "Trucmush",
			CreatedAt: time.Now(),
		},
		&User{
			ID:        "makoush",
			Email:     "makoush@gmail.com",
			FirstName: "Marie",
			LastName:  "Koushtoala",
			CreatedAt: time.Now(),
		})

	assert.Nil(t, err)

	// Count users
	var c int
	c, err = suite.db.UsersCol().Count()
	assert.Nil(t, err)

	assert.Equal(t, c, 2)

	// Fetch one user
	var userJC User
	err = suite.db.UsersCol().Find(bson.M{"first_name": "Jean-Claude"}).One(&userJC)
	assert.Nil(t, err)

	assert.Equal(t, userJC.FirstName, "Jean-Claude")
	assert.Equal(t, userJC.LastName, "Trucmush")
	assert.NotNil(t, userJC.ID)
	assert.NotNil(t, userJC.CreatedAt)

	// Fetch several users
	var allUsers UsersList
	suite.db.UsersCol().Find(nil).All(&allUsers)

	// str, _ := json.MarshalIndent(allUsers, "", " ")
	// fmt.Printf("%s\n", str)

	assert.Equal(t, len(allUsers), 2)
}

func (suite *UserTestSuite) TestJSON() {
	var err error
	var result []byte

	t := suite.T()

	// User
	user := &User{FirstName: "Jean-Claude", LastName: "Trucmush", CreatedAt: time.Now()}

	result, err = json.Marshal(user)
	assert.Nil(t, err)
	assert.NotEmpty(t, result)

	// fmt.Printf("%s\n", result)

	// UsersList
	users := UsersList{
		&User{FirstName: "Jean-Claude", LastName: "Trucmush", CreatedAt: time.Now()},
		&User{FirstName: "Marie", LastName: "Koushtoala", CreatedAt: time.Now()},
	}

	result, err = json.Marshal(users)
	assert.Nil(t, err)
	assert.NotEmpty(t, result)

	// fmt.Printf("%s\n", result)
}
