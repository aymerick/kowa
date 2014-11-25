package models

import (
	"log"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

type DBSession struct {
	mongoSession *mgo.Session
	dbName       string
}

func init() {
	// DEBUG:
	// logout := log.New(os.Stdout, "MGO: ", log.Lshortfile)
	// mgo.SetLogger(logout)
	// mgo.SetDebug(true)
}

// returns a database session
func NewDBSession() *DBSession {
	return &DBSession{
		mongoSession: NewMongoDBSession(),
	}
}

// returns a new mongodb session
func NewMongoDBSession() *mgo.Session {
	// get URI
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = viper.GetString("mongodb_uri")

		if uri == "" {
			log.Fatalln("No connection uri for MongoDB provided")
		}
	}

	return MongoDBSessionForURI(uri)
}

// returns a specific mongodb session
func MongoDBSessionForURI(uri string) *mgo.Session {
	// create session
	result, err := mgo.Dial(uri)
	if result == nil || err != nil {
		log.Fatalf("Can't connect to mongo, go error %v\n", err)
	}

	result.SetSafe(&mgo.Safe{})
	result.SetMode(mgo.Monotonic, true)

	return result
}

//
// DBSession
//

// ensure indexes on all collections
func (this *DBSession) EnsureIndexes() {
	this.EnsureUsersIndexes()
	this.EnsureSitesIndexes()
}

// returns a database handler
func (this *DBSession) DB() *mgo.Database {
	return this.mongoSession.DB(this.DBName())
}

// set database name
func (this *DBSession) SetDBName(name string) {
	this.dbName = name
}

// get database name
func (this *DBSession) DBName() string {
	if this.dbName == "" {
		this.dbName = viper.GetString("mongodb_dbname")
	}

	return this.dbName
}
