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

func (session *DBSession) Copy() *DBSession {
	return &DBSession{
		mongoSession: session.mongoSession.Copy(),
		dbName:       session.dbName,
	}
}

func (session *DBSession) Close() {
	session.mongoSession.Close()
}

// ensure indexes on all collections
func (session *DBSession) EnsureIndexes() {
	session.EnsureUsersIndexes()
	session.EnsureSitesIndexes()
	session.EnsurePostsIndexes()
	session.EnsureEventsIndexes()
	session.EnsurePagesIndexes()
	session.EnsureActivitiesIndexes()
}

// returns a database handler
func (session *DBSession) DB() *mgo.Database {
	return session.mongoSession.DB(session.DBName())
}

// set database name
func (session *DBSession) SetDBName(name string) {
	session.dbName = name
}

// get database name
func (session *DBSession) DBName() string {
	if session.dbName == "" {
		session.dbName = viper.GetString("mongodb_dbname")
	}

	return session.dbName
}
