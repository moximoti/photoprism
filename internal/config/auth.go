package config

import (
	"net/url"
	"path"
	"regexp"

	"github.com/photoprism/photoprism/pkg/rnd"
	"golang.org/x/crypto/bcrypt"
)

func isBcrypt(s string) bool {
	b, err := regexp.MatchString(`^\$2[ayb]\$.{56}$`, s)
	if err != nil {
		return false
	}
	return b
}

// CheckPassword compares given password p with the admin password
func (c *Config) CheckPassword(p string) bool {
	ap := c.AdminPassword()

	if isBcrypt(ap) {
		err := bcrypt.CompareHashAndPassword([]byte(ap), []byte(p))
		return err == nil
	}

	return ap == p
}

// InvalidDownloadToken tests if the token is invalid.
func (c *Config) InvalidDownloadToken(t string) bool {
	return c.DownloadToken() != t
}

// DownloadToken returns the DOWNLOAD api token (you can optionally use a static value for permanent caching).
func (c *Config) DownloadToken() string {
	if c.options.DownloadToken == "" {
		c.options.DownloadToken = rnd.Token(8)
	}

	return c.options.DownloadToken
}

// InvalidPreviewToken tests if the preview token is invalid.
func (c *Config) InvalidPreviewToken(t string) bool {
	return c.PreviewToken() != t && c.DownloadToken() != t
}

// PreviewToken returns the preview image api token (based on the unique storage serial by default).
func (c *Config) PreviewToken() string {
	if c.options.PreviewToken == "" {
		if c.Public() {
			c.options.PreviewToken = "public"
		} else {
			c.options.PreviewToken = c.SerialChecksum()
		}
	}

	return c.options.PreviewToken
}

type AuthConfig struct {
	disableRegistration bool
	passwordPolicy      string
	authProvider        string
	clientID            string
	clientSecret        string
	discoveryEndpoint   string
	callbackUrl         string
}

const (
	RegistrationEnabled  = "enable"
	RegistrationDisabled = "disable"
)

// AuthConfig returns a convenient struct containing relevant authentication settings
func (c *Config) AuthConfig() *AuthConfig {
	u, err := url.Parse(c.Options().SiteUrl)
	if err != nil {
		log.Fatal(err)
	}
	u.Path = path.Join(u.Path, "/api/v1/auth/callback")
	return &AuthConfig{
		disableRegistration: c.Options().DisableRegistration,
		authProvider:        c.Options().AuthProvider,
		clientID:            c.Options().OAuth2ClientID,
		clientSecret:        c.Options().OAuth2ClientSecret,
		discoveryEndpoint:   c.Options().OIDCDiscoveryEndpoint,
		callbackUrl:         u.String(),
	}
}

func (c *AuthConfig) ClientID() string {
	return c.clientID
}
func (c *AuthConfig) ClientSecret() string {
	return c.clientSecret
}
func (c *AuthConfig) AuthProvider() string {
	return c.authProvider
}
func (c *AuthConfig) CallbackUrl() string {
	return c.callbackUrl
}
func (c *AuthConfig) PasswordPolicy() string {
	return c.passwordPolicy
}
func (c *AuthConfig) DiscoveryEndpoint() string {
	return c.discoveryEndpoint
}
