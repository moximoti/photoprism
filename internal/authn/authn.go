package authn

import (
	"errors"
	log "github.com/dsoprea/go-logging"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/service"
	"net/http"
)

const callbackUrl = "http://localhost:2342/api/v1/auth/callback"

func Init() {
	authConf := service.Config().Settings().Auth
	//providerString := authConf.AuthProvider
	err := setGothProvider(authConf)
	if err != nil {
		log.Errorf(err.Error())
	}
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return authConf.AuthProvider, nil
	}
}

func setGothProvider(as config.AuthSettings) error {
	switch as.AuthProvider {
	case config.ProviderOidc:
		// OpenID Connect is based on OpenID Connect Auto Discovery URL (https://openid.net/specs/openid-connect-discovery-1_0-17.html)
		// because the OpenID Connect provider initialize it self in the New(), it can return an error which should be handled or ignored
		// ignore the error for now
		openidConnect, err := openidConnect.New(as.ClientId, as.ClientSecret, callbackUrl, as.DiscoveryEndpoint)
		if openidConnect != nil {
			goth.UseProviders(openidConnect)
		}
		return err
	case config.ProviderGoogle:
		return nil
	case config.ProviderGithub:
		return nil
	default:
		return errors.New("no provider selected")
	}

}

// StartAuthFlow must be called when external authentication is requested
func StartAuthFlow(res http.ResponseWriter, req *http.Request) error {
	url, err := gothic.GetAuthURL(res, req)
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
