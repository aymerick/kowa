package server

import (
	"time"

	"github.com/RangelReale/osin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/aymerick/kowa/models"
)

type oauthStorage struct {
	session *mgo.Session
}

type accessDoc struct {
	ClientID     string    // Client id
	AccessToken  string    // Access token
	RefreshToken string    // Refresh Token. Can be blank
	ExpiresIn    int32     // Token expiration in seconds
	Scope        string    // Requested scope
	CreatedAt    time.Time // Date created
	UserData     interface{}
}

const (
	oauthDBName           = "kowa_oauth"
	clientsCol            = "clients"
	authorizationsColName = "authorizations" // NOT USED
	accessesColName       = "accesses"
)

const (
	refreshToken      = "refreshtoken"
	oauthClientID     = "kowa"
	oauthClientSecret = "none"
)

func newOAuthStorage() *oauthStorage {
	storage := &oauthStorage{
		session: models.NewMongoDBSession(),
	}

	index := mgo.Index{
		Key:        []string{refreshToken},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}

	accesses := storage.session.DB(oauthDBName).C(accessesColName)
	err := accesses.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	return storage
}

func (storage *oauthStorage) EnsureOAuthClient() error {
	client := &oauthClient{
		ID:          oauthClientID,
		Secret:      oauthClientSecret,
		RedirectURI: "http://localhost:35830/", // @todo Check that
	}

	_, err := storage.clientsCol().UpsertId(client.ID, client)
	return err
}

func (storage *oauthStorage) DB() *mgo.Database {
	return storage.session.DB(oauthDBName)
}

func (storage *oauthStorage) clientsCol() *mgo.Collection {
	return storage.DB().C(clientsCol)
}

func (storage *oauthStorage) authorizationsCol() *mgo.Collection {
	return storage.DB().C(authorizationsColName)
}

func (storage *oauthStorage) accessesCol() *mgo.Collection {
	return storage.DB().C(accessesColName)
}

//
// Implements osin.Storage interface
//

func (storage *oauthStorage) Clone() osin.Storage {
	return &oauthStorage{session: storage.session.Copy()}
}

func (storage *oauthStorage) Close() {
	storage.session.Close()
}

func (storage *oauthStorage) GetClient(id string) (osin.Client, error) {
	client := &oauthClient{}
	err := storage.clientsCol().FindId(id).One(client)
	return client, err
}

// NOT USED
func (storage *oauthStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	_, err := storage.authorizationsCol().UpsertId(data.Code, data)
	return err
}

// NOT USED
func (storage *oauthStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	authData := &osin.AuthorizeData{}
	err := storage.authorizationsCol().FindId(code).One(authData)
	return authData, err
}

// NOT USED
func (storage *oauthStorage) RemoveAuthorize(code string) error {
	return storage.authorizationsCol().RemoveId(code)
}

func (storage *oauthStorage) SaveAccess(data *osin.AccessData) error {
	doc := &accessDoc{
		ClientID:     data.Client.GetId(),
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		ExpiresIn:    data.ExpiresIn,
		Scope:        data.Scope,
		CreatedAt:    data.CreatedAt,
		UserData:     data.UserData,
	}
	_, err := storage.accessesCol().UpsertId(data.AccessToken, doc)
	return err
}

func (storage *oauthStorage) docToAccessData(doc *accessDoc) (*osin.AccessData, error) {
	var result *osin.AccessData
	var err error

	var client osin.Client
	client, err = storage.GetClient(doc.ClientID)
	if err == nil {
		result = &osin.AccessData{
			Client:       client,
			AccessToken:  doc.AccessToken,
			RefreshToken: doc.RefreshToken,
			ExpiresIn:    doc.ExpiresIn,
			Scope:        doc.Scope,
			CreatedAt:    doc.CreatedAt,
			UserData:     doc.UserData,
		}
	}

	return result, err
}

func (storage *oauthStorage) LoadAccess(token string) (*osin.AccessData, error) {
	var result *osin.AccessData
	var err error

	doc := &accessDoc{}
	err = storage.accessesCol().FindId(token).One(doc)
	if err == nil {
		result, err = storage.docToAccessData(doc)
	}

	return result, err
}

func (storage *oauthStorage) RemoveAccess(token string) error {
	return storage.accessesCol().RemoveId(token)
}

func (storage *oauthStorage) LoadRefresh(token string) (*osin.AccessData, error) {
	var result *osin.AccessData
	var err error

	doc := &accessDoc{}
	err = storage.accessesCol().Find(bson.M{refreshToken: token}).One(doc)
	if err == nil {
		result, err = storage.docToAccessData(doc)
	}

	return result, err
}

func (storage *oauthStorage) RemoveRefresh(token string) error {
	selector := bson.M{refreshToken: token}
	modifier := bson.M{"$unset": bson.M{refreshToken: 1}}

	return storage.accessesCol().Update(selector, modifier)
}
