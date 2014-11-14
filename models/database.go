package models

import (
	"log"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

var (
	dbSession *mgo.Session
	dbName    string
)

func init() {
	// @todo EnsureIndexes()
}

func MongoDBSessionForURI(uri string) *mgo.Session {
	// create session
	session, err := mgo.Dial(uri)
	if session == nil || err != nil {
		log.Fatalf("Can't connect to mongo, go error %v\n", err)
	}

	session.SetSafe(&mgo.Safe{})
	session.SetMode(mgo.Monotonic, true)

	return session
}

func SetDBSession(session *mgo.Session) {
	dbSession = session
}

func DBSession() *mgo.Session {
	if dbSession == nil {
		// get URI
		uri := os.Getenv("MONGODB_URI")
		if uri == "" {
			uri = viper.GetString("mongodb_uri")

			if uri == "" {
				log.Fatalln("No connection uri for MongoDB provided")
			}
		}

		SetDBSession(MongoDBSessionForURI(uri))
	}

	return dbSession
}

func SetDBName(name string) {
	dbName = name
}

func DBName() string {
	if dbName == "" {
		dbName = viper.GetString("mongodb_dbname")
	}

	return dbName
}

func DB() *mgo.Database {
	return DBSession().DB(DBName())
}
