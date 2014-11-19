package server

import (
	"github.com/RangelReale/osin"
	"github.com/aymerick/kowa/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type OAuthStorage struct {
	session *mgo.Session
}

const (
	DBNAME             = "kowa_oauth"
	CLIENTS_COL        = "clients"
	AUTHORIZATIONS_COL = "authorizations"
	ACCESSES_COL       = "accesses"
)

const (
	REFRESHTOKEN      = "refreshtoken"
	CLIENT_ID         = "kowa"
	CLIENT_SECRET     = "none"
	CLIENT_AUTH_VALUE = "a293YTpub25l" // This the base64 value of <CLIENT_ID>:<CLIENT_SECRET>
)

func NewOAuthStorage() *OAuthStorage {
	storage := &OAuthStorage{session: models.DBSession()}

	index := mgo.Index{
		Key:        []string{REFRESHTOKEN},
		Unique:     false, // refreshtoken is sometimes empty
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}

	accesses := storage.session.DB(DBNAME).C(ACCESSES_COL)
	err := accesses.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	return storage
}

func (this *OAuthStorage) DB() *mgo.Database {
	return this.session.DB(DBNAME)
}

func (this *OAuthStorage) clientsCol() *mgo.Collection {
	return this.DB().C(CLIENTS_COL)
}

func (this *OAuthStorage) authorizationsCol() *mgo.Collection {
	return this.DB().C(AUTHORIZATIONS_COL)
}

func (this *OAuthStorage) accessesCol() *mgo.Collection {
	return this.DB().C(ACCESSES_COL)
}

func (this *OAuthStorage) SetClient(id string, client osin.Client) error {
	_, err := this.clientsCol().UpsertId(id, client)
	return err
}

func (this *OAuthStorage) SetupDefaultClient() (osin.Client, error) {
	client := &osin.DefaultClient{
		Id:          CLIENT_ID,
		Secret:      CLIENT_SECRET,
		RedirectUri: "http://localhost:35830/appauth", // @todo Check that
	}
	err := this.SetClient("kowa", client)

	return client, err
}

//
// Implements osin.Storage interface
//

func (this *OAuthStorage) Clone() osin.Storage {
	return &OAuthStorage{session: this.session.Copy()}
}

func (this *OAuthStorage) Close() {
	this.session.Close()
}

func (this *OAuthStorage) GetClient(id string) (osin.Client, error) {
	client := &osin.DefaultClient{}
	err := this.clientsCol().FindId(id).One(client)
	return client, err
}

func (this *OAuthStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	_, err := this.authorizationsCol().UpsertId(data.Code, data)
	return err
}

func (this *OAuthStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	authData := &osin.AuthorizeData{}
	err := this.authorizationsCol().FindId(code).One(authData)
	return authData, err
}

func (this *OAuthStorage) RemoveAuthorize(code string) error {
	return this.authorizationsCol().RemoveId(code)
}

func (this *OAuthStorage) SaveAccess(data *osin.AccessData) error {
	_, err := this.accessesCol().UpsertId(data.AccessToken, data)
	return err
}

func (this *OAuthStorage) LoadAccess(token string) (*osin.AccessData, error) {
	accData := &osin.AccessData{}
	err := this.accessesCol().FindId(token).One(accData)
	return accData, err
}

func (this *OAuthStorage) RemoveAccess(token string) error {
	return this.accessesCol().RemoveId(token)
}

func (this *OAuthStorage) LoadRefresh(token string) (*osin.AccessData, error) {
	accData := &osin.AccessData{}
	err := this.accessesCol().Find(bson.M{REFRESHTOKEN: token}).One(accData)
	return accData, err
}

func (this *OAuthStorage) RemoveRefresh(token string) error {
	selector := bson.M{REFRESHTOKEN: token}
	modifier := bson.M{"$unset": bson.M{REFRESHTOKEN: 1}}

	return this.accessesCol().Update(selector, modifier)
}
