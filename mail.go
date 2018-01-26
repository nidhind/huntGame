package main

import (
	"errors"
	"html/template"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailObj struct {
	To          string
	Cc          string
	Bcc         string
	Subject     string
	MessageText string
	MessageHTML string
}

var passwordResetTemplate *template.Template

func Send(m *MailObj) error {
	// TODO: Move 'from' parameters to config
	from := mail.NewEmail("Quest Admin", "questmaster2018@gmail.com")

	to := mail.NewEmail("", m.To)
	subject := m.Subject
	content := mail.NewContent("text/html", m.MessageHTML)

	message := mail.NewV3MailInit(from, subject, to, content)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)

	if err != nil {
		return err
	}
	if response.StatusCode/100 != 2 {
		return errors.New(response.Body)
	}

	return nil
}

func getResetPswdEmailTemplate() (*template.Template, error) {
	var pt *template.Template
	var err error

	if passwordResetTemplate != pt {
		return passwordResetTemplate, nil
	}

	path := os.Getenv("RESET_PSWD_TEMPLATE")
	if path != "" {
		pt, err = template.ParseFiles(path)
	} else {
		t := `<!DOCTYPE html>
<html>
<head>
	<title>Forgot Password</title>
</head>
<body>
	<h1>Hi {{.Email}}</h1>
	<h3>Forgot your password...??</h3>
	<h3>No worries...Click <a href="{{.ResetLink}}">here</a> to reset your password</h3>
	<h3>Trouble with the link? Please copy and paste the following in a browser</h3>
	<p>{{.ResetLink}}</p>
</body>
</html>
`
		pt, err = template.New("pswd_reset_mail").Parse(t)
	}

	if err != nil {
		return pt, err
	}
	passwordResetTemplate = pt
	return passwordResetTemplate, nil

}
