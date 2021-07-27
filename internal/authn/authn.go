package authn

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/photoprism/photoprism/internal/event"
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

var providerConfig authnConfig

func Init(authConf authnConfig) error {
	providerConfig = authConf
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

func getJwkEndpoint() (string, error) {
	res, err := http.Get(providerConfig.DiscoveryEndpoint())
	if err != nil {
		return "", err
	}

	wellknown := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&wellknown)
	if err != nil {
		return "", err
	}

	jwkurl, ok := wellknown["jwks_uri"].(string)
	if !ok {
		return "", errors.New("couldn't retrieve public key to verify id_token")
	}
	return jwkurl, nil
}

// validates ID token and extracts external user ID extracted from ID token
func ValidateAndExtractID(idToken string) (string, error) {
	// evaluate https://github.com/MicahParks/keyfunc
	// referenced here https://stackoverflow.com/questions/41077953/go-language-and-verify-jwt
	// only support asymmetric crypto signature to mitigate attack vector
	token, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		event.Log.Debugf("Token Signing method: %s", token.Method)
		event.Log.Debugf("Token Header: %s", token.Header)

		keyID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("expecting JWT header to have string kid")
		}

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("signing method not supported")
		}

		jwkurl, err := getJwkEndpoint()
		if err != nil {
			return nil, err
		}

		keyset, err := jwk.Fetch(context.Background(), jwkurl)
		if err != nil {
			return nil, err
		}

		if key, ok := keyset.LookupKeyID(keyID); ok {
			rsakey := new(rsa.PublicKey)
			err := key.Raw(rsakey)
			if err != nil {
				return "", err
			}

			return rsakey, nil
		}
		return nil, errors.New("couldn't determine key")
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		}
	}

	return "", errors.New("token invalid")
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
