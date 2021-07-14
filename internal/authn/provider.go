package authn

// Provider needs to be implemented for each 3rd party authentication provider
// e.g. Facebook, Twitter, etc...
type Provider interface {
	AuthCodeURL(state string) (string, error)
	BeginCodeFlow(state string) (string, error)
	//Name() string
	//SetName(name string)
	//BeginAuth(state string) (Session, error)
	//UnmarshalSession(string) (Session, error)
	//GetUser() (entity.User, error)
	//Debug(bool)
	//RefreshToken(refreshToken string) (*oauth2.Token, error) //Get new access token based on the refresh token
	//RefreshTokenAvailable() bool                             //Refresh token is provided by auth provider or not
}
