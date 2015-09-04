package server

type oauthClient struct {
	ID          string      `bson:"_id"`
	Secret      string      `bson:"secret"`
	RedirectURI string      `bson:"redirect_uri"`
	UserData    interface{} `bson:"user_data"`
}

func (client *oauthClient) GetId() string {
	return client.ID
}

func (client *oauthClient) GetSecret() string {
	return client.Secret
}

func (client *oauthClient) GetRedirectUri() string {
	return client.RedirectURI
}

func (client *oauthClient) GetUserData() interface{} {
	return client.UserData
}
