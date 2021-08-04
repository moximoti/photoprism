package mail

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/event"
	"html/template"
	"net"
	"net/smtp"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

type Config interface {
	MailHost() string
	MailFrom() string
	MailUsername() string
	MailPassword() string
	MailTemplatesPath() string
	SiteUrl() string
}

var conf Config
var log = event.Log
var templates *template.Template

func loadTemplates() error {
	paths := []string{
		filepath.Join(conf.MailTemplatesPath(), "reset_password.tmpl"),
		filepath.Join(conf.MailTemplatesPath(), "confirm_user.tmpl"),
	}

	t, err := template.ParseFiles(paths...)
	if err != nil {
		return err
	}
	templates = t
	return nil
}

func getUrl(relpath string) string {
	u, _ := url.Parse(conf.SiteUrl())
	u.Path = path.Join(u.Path, relpath)
	return u.String()
}

func Init(c Config) error {
	log.Debugf("Init Mailer")
	conf = c
	err := loadTemplates()
	if err != nil {
		return err
	}
	//sendTestMail()
	return checkConnection()
}

func sendTestMail() error {
	testUser := entity.User{
		UserName:     "goofy",
		FullName:     "Timo Volkes",
		PrimaryEmail: "timo.tvm@gmail.com",
	}
	return SendUserConfirmationMail(&testUser)
}

func SendMail(to []string, msg []byte) error {
	client, err := connect()
	if err != nil {
		return err
	}
	defer client.Close()

	for _, rec := range to {
		if err = client.Rcpt(rec); err != nil {
			return fmt.Errorf("Rcpt: %v", err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("Data: %v", err)
	}
	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("Write: %v", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("Close: %v", err)
	}

	return nil
}

func connect() (*smtp.Client, error) {
	if conf == nil {
		return nil, errors.New("mailer is not configured")
	}
	hostPort := strings.Split(conf.MailHost(), ":")

	tlsconfig := &tls.Config{
		InsecureSkipVerify: false, // TODO offer config option skipVerifyTLS
		ServerName:         hostPort[0],
	}

	conn, err := net.Dial("tcp", conf.MailHost())
	if err != nil {
		return nil, err
	}
	//defer conn.Close()

	isSecureConn := hostPort[1] == "465" // TODO offer config option forceTLS
	// Start TLS directly if the port ends with 465 (SMTPS protocol)
	if isSecureConn {
		conn = tls.Client(conn, tlsconfig)
	}

	client, err := smtp.NewClient(conn, conf.MailHost())
	if err != nil {
		return nil, fmt.Errorf("NewClient: %v", err)
	}

	// If not using SMTPS, always use STARTTLS if available
	hasStartTLS, _ := client.Extension("STARTTLS")
	if !isSecureConn && hasStartTLS {
		if err = client.StartTLS(tlsconfig); err != nil {
			return nil, fmt.Errorf("StartTLS: %v", err)
		}
	}

	canAuth, options := client.Extension("AUTH")
	if canAuth {
		var auth smtp.Auth
		if strings.Contains(options, "CRAM-MD5") {
			auth = smtp.CRAMMD5Auth(conf.MailUsername(), conf.MailPassword())
		} else if strings.Contains(options, "PLAIN") {
			auth = smtp.PlainAuth("", conf.MailUsername(), conf.MailPassword(), hostPort[0])
		}
		if auth != nil {
			if err = client.Auth(auth); err != nil {
				return nil, fmt.Errorf("Auth: %v", err)
			}
		}
	}

	if err = client.Mail(conf.MailFrom()); err != nil {
		return nil, fmt.Errorf("Mail: %v", err)
	}
	return client, nil
}

func checkConnection() error {
	client, err := connect()
	if err != nil {
		return err
	}
	defer client.Close()
	return client.Noop()
}
