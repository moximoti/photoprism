package authn

import (
	"errors"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/photoprism/photoprism/internal/authn/provider"
	"github.com/photoprism/photoprism/internal/service"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

const callbackUrl = "http://localhost:2342/api/v1/auth/callback"

//var (
//	OauthConfig = &oauth2.Config{
//		RedirectURL:    callbackUrl,
//		ClientID:     "photoprism-dev",
//		ClientSecret: "341e8af4-4ed7-40cc-bd2c-b29a5e0cd40c",
//		Scopes:       []string{"profile", "email", "openid"},
//		Endpoint:     oauth2.Endpoint{
//			AuthURL:   "https://keycloak.timovolkmann.de/auth/realms/master/protocol/openid-connect/auth",
//			TokenURL:  "https://keycloak.timovolkmann.de/auth/realms/master/protocol/openid-connect/token",
//			AuthStyle: 0,
//		},
//	}
//	// Some random string, random for each request
//	OauthStateString = "random"
//)
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

var authProvider Provider

func Init() {
	authConf := service.Config().Settings().Auth
	//providerString := authConf.AuthProvider
	authProvider, _ = provider.NewOidc(authConf, callbackUrl)
}

func StartAuthFlow(res http.ResponseWriter, req *http.Request) (string, error) {
	state := "randomsstate"
	// TODO: randomize state parameter and implement some kind of session to glue StartAuthFlow and FinalizeAuth together

	if url, err := authProvider.AuthCodeURL(state); err != nil {
		return "", err
	} else {
		return url, nil
	}
}

func FinalizeAuthFlow(res http.ResponseWriter, req *http.Request) (ExternalUser, error) {

}

type FlowStateCache map[string]FlowState

var flowStore FlowStateCache

func (c FlowStateCache) Set(key string, value FlowState) error {
	c[key] = value
	return nil
}

func (c FlowStateCache) Get(key string) (FlowState, error) {
	if res, ok := c[key]; ok {
		return res, nil
	}
	return nil, errors.New("no cached flow state")
}
