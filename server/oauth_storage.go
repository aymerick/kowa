package server

import (
	"time"

	"github.com/RangelReale/osin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/aymerick/kowa/models"
)

type OAuthStorage struct {
	session *mgo.Session
}

type AccessDoc struct {
	ClientId     string    // Client id
	AccessToken  string    // Access token
	RefreshToken string    // Refresh Token. Can be blank
	ExpiresIn    int32     // Token expiration in seconds
	Scope        string    // Requested scope
	CreatedAt    time.Time // Date created
	UserData     interface{}
}

const (
	OAUTH_DBNAME       = "kowa_oauth"
	CLIENTS_COL        = "clients"
	AUTHORIZATIONS_COL = "authorizations" // NOT USED
	ACCESSES_COL       = "accesses"
)

const (
	REFRESHTOKEN        = "refreshtoken"
	OAUTH_CLIENT_ID     = "kowa"
	OAUTH_CLIENT_SECRET = "none"
)

func NewOAuthStorage() *OAuthStorage {
	storage := &OAuthStorage{
		session: models.NewMongoDBSession(),
	}

	index := mgo.Index{
		Key:        []string{REFRESHTOKEN},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}

	accesses := storage.session.DB(OAUTH_DBNAME).C(ACCESSES_COL)
	err := accesses.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	return storage
}

func (storage *OAuthStorage) EnsureOAuthClient() error {
	client := &OAuthClient{
		Id:          OAUTH_CLIENT_ID,
		Secret:      OAUTH_CLIENT_SECRET,
		RedirectUri: "http://localhost:35830/", // @todo Check that
	}

	_, err := storage.clientsCol().UpsertId(client.Id, client)
	return err
}

func (storage *OAuthStorage) DB() *mgo.Database {
	return storage.session.DB(OAUTH_DBNAME)
}

func (storage *OAuthStorage) clientsCol() *mgo.Collection {
	return storage.DB().C(CLIENTS_COL)
}

func (storage *OAuthStorage) authorizationsCol() *mgo.Collection {
	return storage.DB().C(AUTHORIZATIONS_COL)
}

func (storage *OAuthStorage) accessesCol() *mgo.Collection {
	return storage.DB().C(ACCESSES_COL)
}

//
// Implements osin.Storage interface
//

func (storage *OAuthStorage) Clone() osin.Storage {
	return &OAuthStorage{session: storage.session.Copy()}
}

func (storage *OAuthStorage) Close() {
	storage.session.Close()
}

func (storage *OAuthStorage) GetClient(id string) (osin.Client, error) {
	client := &OAuthClient{}
	err := storage.clientsCol().FindId(id).One(client)
	return client, err
}

// NOT USED
func (storage *OAuthStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	_, err := storage.authorizationsCol().UpsertId(data.Code, data)
	return err
}

// NOT USED
func (storage *OAuthStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	authData := &osin.AuthorizeData{}
	err := storage.authorizationsCol().FindId(code).One(authData)
	return authData, err
}

// NOT USED
func (storage *OAuthStorage) RemoveAuthorize(code string) error {
	return storage.authorizationsCol().RemoveId(code)
}

func (storage *OAuthStorage) SaveAccess(data *osin.AccessData) error {
	doc := &AccessDoc{
		ClientId:     data.Client.GetId(),
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

func (storage *OAuthStorage) docToAccessData(doc *AccessDoc) (*osin.AccessData, error) {
	var result *osin.AccessData
	var err error

	var client osin.Client
	client, err = storage.GetClient(doc.ClientId)
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

func (storage *OAuthStorage) LoadAccess(token string) (*osin.AccessData, error) {
	var result *osin.AccessData
	var err error

	doc := &AccessDoc{}
	err = storage.accessesCol().FindId(token).One(doc)
	if err == nil {
		result, err = storage.docToAccessData(doc)
	}

	return result, err
}

func (storage *OAuthStorage) RemoveAccess(token string) error {
	return storage.accessesCol().RemoveId(token)
}

func (storage *OAuthStorage) LoadRefresh(token string) (*osin.AccessData, error) {
	var result *osin.AccessData
	var err error

	doc := &AccessDoc{}
	err = storage.accessesCol().Find(bson.M{REFRESHTOKEN: token}).One(doc)
	if err == nil {
		result, err = storage.docToAccessData(doc)
	}

	return result, err
}

func (storage *OAuthStorage) RemoveRefresh(token string) error {
	selector := bson.M{REFRESHTOKEN: token}
	modifier := bson.M{"$unset": bson.M{REFRESHTOKEN: 1}}

	return storage.accessesCol().Update(selector, modifier)
}
