package config

import (
	"github.com/photoprism/photoprism/internal/mail"
	"path/filepath"
)

func (c *Config) MailHost() string {
	return c.Options().MailHost
}

func (c *Config) MailUsername() string {
	return c.Options().MailUsername
}

func (c *Config) MailPassword() string {
	return c.Options().MailPassword
}

func (c *Config) MailFrom() string {
	return c.Options().MailFrom
}

// MailTemplatesPath returns the mail templates path.
func (c *Config) MailTemplatesPath() string {
	return filepath.Join(c.AssetsPath(), "mail")
}

func (c *Config) initMail() {
	err := mail.Init(c)
	if err != nil {
		if c.EmailConfirmationEnabled() {
			panic(err)
		} else {
			log.Warning(err)
		}
	}
}
