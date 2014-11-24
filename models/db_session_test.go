package models

const (
	TEST_MONGODB_URI = "mongodb://localhost:27017/"
	TEST_DBNAME      = "kowa_test"
)

func NewTestDBSession() *DBSession {
	return &DBSession{mongoSession: MongoDBSessionForURI(TEST_MONGODB_URI)}
}
