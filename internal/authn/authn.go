package authn

import (
	"errors"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/photoprism/photoprism/pkg/rnd"
	"github.com/sirupsen/logrus"
	"net/http"
)

type authnConfig interface {
	ClientID() string
	ClientSecret() string
	AuthProvider() string
	CallbackUrl() string
	DiscoveryEndpoint() string
}

func Init(authConf authnConfig) error {
	//authConf := service.Config().AuthConfig()
	err := setGothProvider(authConf)
	if err != nil {
		return err
	}
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return authConf.AuthProvider(), nil
	}
	gothic.Store = sessions.NewCookieStore([]byte(rnd.UUID()))
	return nil
}

const (
	ProviderNone   = "none"
	ProviderOidc   = "openid-connect"
	ProviderGoogle = "google"
	ProviderGithub = "github"
)

func setGothProvider(ac authnConfig) error {
	switch ac.AuthProvider() {
	case ProviderOidc:
		// OpenID Connect is based on OpenID Connect Auto Discovery URL (https://openid.net/specs/openid-connect-discovery-1_0-17.html)
		// because the OpenID Connect provider initialize it self in the New(), it can return an error which should be handled or ignored
		// ignore the error for now
		oidc, err := openidConnect.New(ac.ClientID(), ac.ClientSecret(), ac.CallbackUrl(), ac.DiscoveryEndpoint())
		if oidc != nil {
			goth.UseProviders(oidc)
		}
		return err
	case ProviderGoogle:
		return nil
	case ProviderGithub:
		return nil
	default:
		return errors.New("no provider selected")
	}

}

// StartAuthFlow must be called when external authentication is requested
func StartAuthFlow(res http.ResponseWriter, req *http.Request) error {
	url, err := gothic.GetAuthURL(res, req)
	logrus.Debugf("Auth URL: %s", url)
	if err == nil {
		http.Redirect(res, req, url, http.StatusTemporaryRedirect)
	}
	return err
}

// FinalizeAuthFlow must be called on provider callback
func FinalizeAuthFlow(res http.ResponseWriter, req *http.Request) (goth.User, error) {
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		return user, err
	}

	return user, nil
}
