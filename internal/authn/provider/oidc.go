package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/photoprism/photoprism/internal/config"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

type OidcProvider struct {
	ClientId     string
	ClientSecret string
	CallbackURL  string
	DiscoveryURL string
	HTTPClient   *http.Client
	config       *oauth2.Config
	openIDConf   *OpenIDConfig
	//providerName string
	//profileURL   string
	//emailURL     string
}

type OidcFlowState struct {
}

type OpenIDConfig struct {
	AuthEndpoint     string `json:"authorization_endpoint"`
	TokenEndpoint    string `json:"token_endpoint"`
	UserInfoEndpoint string `json:"userinfo_endpoint"`

	// If OpenID discovery is enabled, the end_session_endpoint field can optionally be provided
	// in the discovery endpoint response according to OpenID spec. See:
	// https://openid.net/specs/openid-connect-session-1_0-17.html#OPMetadata
	EndSessionEndpoint string `json:"end_session_endpoint,omitempty"`
	Issuer             string `json:"issuer"`
}

func NewOidc(settings config.AuthSettings, callbackUrl string) (*OidcProvider, error) {
	provider := &OidcProvider{
		HTTPClient: http.DefaultClient,
	}
	if settings.DiscoveryEndpoint != "" {
		provider.DiscoveryURL = settings.DiscoveryEndpoint
		openidconf, err := provider.fetchOidcConfiguration()
		if err != nil {
			return nil, err
		}
		provider.openIDConf = openidconf
		provider.config = &oauth2.Config{
			ClientID:     settings.ClientId,
			ClientSecret: settings.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:   openidconf.AuthEndpoint,
				TokenURL:  openidconf.TokenEndpoint,
				AuthStyle: 0,
			},
			RedirectURL: callbackUrl,
			Scopes:      []string{"profile", "email", "openid"},
		}
	}
	return provider, nil
}
func (p *OidcProvider) AuthCodeURL(state string) (string, error) {
	return p.config.AuthCodeURL(state), nil
}

type UserInfo struct {
	UID   string
	email string
}

func (p *OidcProvider) FinalizeAuthFlow(state string, code string) (*UserInfo, error) {
	token, err := p.config.Exchange(context.Background(), code)
	if err != nil {
		return UserInfo{}, err
	}
	return p.resolveUser(token)
}

func (p *OidcProvider) resolveUser(token *oauth2.Token) (*UserInfo, error) {
	if p.openIDConf.UserInfoEndpoint == "" {
		return nil, errors.New("No UserInfo Endpoint available")
		// TODO: alternatively extract information from ID token
	}
	url := p.openIDConf.UserInfoEndpoint
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &UserInfo{}, nil
}

func (p *OidcProvider) fetchOidcConfiguration() (*OpenIDConfig, error) {
	res, err := p.HTTPClient.Get(p.DiscoveryURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	openIDConfig := &OpenIDConfig{}
	err = json.Unmarshal(body, openIDConfig)
	if err != nil {
		return nil, err
	}

	return openIDConfig, nil
}
