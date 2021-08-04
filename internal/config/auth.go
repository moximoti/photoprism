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

// Satisfying config interfaces for relevant authentication settings

func (c *Config) ClientID() string {
	return c.Options().OAuth2ClientID
}
func (c *Config) ClientSecret() string {
	return c.Options().OAuth2ClientSecret
}
func (c *Config) AuthProvider() string {
	return c.Options().AuthProvider
}
func (c *Config) CallbackUrl() string {
	u, err := url.Parse(c.SiteUrl())
	if err != nil {
		log.Fatal(err)
	}
	u.Path = path.Join(u.Path, "/api/v1/auth/callback")
	return u.String()
}
func (c *Config) PasswordPolicy() string {
	panic("not implemented")
}
func (c *Config) DiscoveryEndpoint() string {
	return c.Options().OIDCDiscoveryEndpoint
}

func (c *Config) RegistrationDisabled() bool {
	return c.Options().DisableRegistration
}
func (c *Config) AdminConfirmationEnabled() bool {
	return c.Options().AdminConfirmation
}
func (c *Config) EmailConfirmationEnabled() bool {
	return c.Options().EmailConfirmation
}
