package mail

import (
	"bytes"
	"fmt"
	"github.com/jordan-wright/email"
	"jaytaylor.com/html2text"
)

type User interface {
	Name() string
	Username() string
	Email() string
	GenerateEmailActivateCode() string
}

func SendUserConfirmationMail(u User) error {
	var bufHtml bytes.Buffer

	data := map[string]string{
		"FullName": u.Name(),
		"Link": fmt.Sprintf("%s?code=%s&username=%s",
			getUrl("/api/v1/auth/activate_user"),
			u.GenerateEmailActivateCode(),
			u.Username()),
	}

	err := templates.ExecuteTemplate(&bufHtml, "confirm_user.tmpl", data)
	if err != nil {
		return err
	}
	text, err := html2text.FromReader(&bufHtml)
	if err != nil {
		return err
	}
	to := []string{u.Email()}
	e := email.NewEmail()
	e.From = fmt.Sprintf("Photoprism <%s>", conf.MailFrom())
	e.To = to
	e.Subject = "E-Mail Confirmation"
	e.Text = []byte(text)
	e.HTML = bufHtml.Bytes()

	msg, err := e.Bytes()
	if err != nil {
		return err
	}

	return SendMail(to, msg)
}

func SendPasswordResetMail(u User) error {
	var bufHtml bytes.Buffer

	data := map[string]string{
		"FullName": u.Name(),
		"Link": fmt.Sprintf("%s?code=%s&username=%s",
			getUrl("/api/v1/auth/reset_password"),
			u.GenerateEmailActivateCode(),
			u.Username()),
	}

	err := templates.ExecuteTemplate(&bufHtml, "reset_password.tmpl", data)
	if err != nil {
		return err
	}
	text, err := html2text.FromReader(&bufHtml)
	if err != nil {
		return err
	}
	to := []string{u.Email()}
	e := email.NewEmail()
	e.From = fmt.Sprintf("Photoprism <%s>", conf.MailFrom())
	e.To = to
	e.Subject = "Password Reset"
	e.Text = []byte(text)
	e.HTML = bufHtml.Bytes()

	msg, err := e.Bytes()
	if err != nil {
		return err
	}

	return SendMail(to, msg)
}
