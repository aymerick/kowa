package server

import (
	"time"

	"github.com/RangelReale/osin"
	"github.com/aymerick/kowa/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	DBNAME             = "kowa_oauth"
	CLIENTS_COL        = "clients"
	AUTHORIZATIONS_COL = "authorizations" // NOT USED
	ACCESSES_COL       = "accesses"
)

const (
	REFRESHTOKEN  = "refreshtoken"
	CLIENT_ID     = "kowa"
	CLIENT_SECRET = "none"
)

func NewOAuthStorage() *OAuthStorage {
	storage := &OAuthStorage{session: models.DBSession().Clone()}

	index := mgo.Index{
		Key:        []string{REFRESHTOKEN},
		Unique:     false,
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
		RedirectUri: "http://localhost:35830/", // @todo Check that
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

// NOT USED
func (this *OAuthStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	_, err := this.authorizationsCol().UpsertId(data.Code, data)
	return err
}

// NOT USED
func (this *OAuthStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	authData := &osin.AuthorizeData{}
	err := this.authorizationsCol().FindId(code).One(authData)
	return authData, err
}

// NOT USED
func (this *OAuthStorage) RemoveAuthorize(code string) error {
	return this.authorizationsCol().RemoveId(code)
}

func (this *OAuthStorage) SaveAccess(data *osin.AccessData) error {
	doc := &AccessDoc{
		ClientId:     data.Client.GetId(),
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		ExpiresIn:    data.ExpiresIn,
		Scope:        data.Scope,
		CreatedAt:    data.CreatedAt,
		UserData:     data.UserData,
	}
	_, err := this.accessesCol().UpsertId(data.AccessToken, doc)
	return err
}

func (this *OAuthStorage) docToAccessData(doc *AccessDoc) (*osin.AccessData, error) {
	var result *osin.AccessData
	var err error

	var client osin.Client
	client, err = this.GetClient(doc.ClientId)
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

func (this *OAuthStorage) LoadAccess(token string) (*osin.AccessData, error) {
	var result *osin.AccessData
	var err error

	doc := &AccessDoc{}
	err = this.accessesCol().FindId(token).One(doc)
	if err == nil {
		result, err = this.docToAccessData(doc)
	}

	return result, err
}

func (this *OAuthStorage) RemoveAccess(token string) error {
	return this.accessesCol().RemoveId(token)
}

func (this *OAuthStorage) LoadRefresh(token string) (*osin.AccessData, error) {
	var result *osin.AccessData
	var err error

	doc := &AccessDoc{}
	err = this.accessesCol().Find(bson.M{REFRESHTOKEN: token}).One(doc)
	if err == nil {
		result, err = this.docToAccessData(doc)
	}

	return result, err
}

func (this *OAuthStorage) RemoveRefresh(token string) error {
	selector := bson.M{REFRESHTOKEN: token}
	modifier := bson.M{"$unset": bson.M{REFRESHTOKEN: 1}}

	return this.accessesCol().Update(selector, modifier)
}
