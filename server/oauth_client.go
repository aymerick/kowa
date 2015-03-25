package server

type OAuthClient struct {
	Id          string      `bson:"_id"`
	Secret      string      `bson:"secret"`
	RedirectUri string      `bson:"redirect_uri"`
	UserData    interface{} `bson:"user_data"`
}

func (client *OAuthClient) GetId() string {
	return client.Id
}

func (client *OAuthClient) GetSecret() string {
	return client.Secret
}

func (client *OAuthClient) GetRedirectUri() string {
	return client.RedirectUri
}

func (client *OAuthClient) GetUserData() interface{} {
	return client.UserData
}
