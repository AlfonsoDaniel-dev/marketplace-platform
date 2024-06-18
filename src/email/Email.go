package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/mail"
	"os"
	"shopperia/src/common/models"
	"strings"
)

type Email struct {
	Headers map[string]string
	Subject string
	From    from
	Dest    dest
	Message string
}

type from struct {
	Name    string
	Email   string
	Address mail.Address
}

type dest struct {
	Name    string
	Email   string
	Address mail.Address
}

func setHeaders(from, to mail.Address, subject string) map[string]string {
	Headers := make(map[string]string)
	Headers["FROM"] = from.String()
	Headers["To"] = to.String()
	Headers["Subject"] = subject
	Headers["Content-Type"] = "text/html; charset=utf-8"

	return Headers
}

func newFrom() from {
	name := strings.TrimSpace(os.Getenv("EMAIL_PROVIDER_NAME"))
	email := strings.TrimSpace(os.Getenv("EMAIL_PROVIDER_ACCOUNT"))

	Fromaddress := mail.Address{Name: name, Address: email}

	return from{
		Name:    name,
		Email:   email,
		Address: Fromaddress,
	}
}

func newDest(name, email string) dest {

	destAddress := mail.Address{Name: name, Address: email}

	return dest{
		Name:    name,
		Email:   email,
		Address: destAddress,
	}
}

func (Email *Email) setSubject(subject string) {
	Email.Subject = subject
}

func setAndParseTemplate(templatePath string) (*template.Template, error) {
	temp, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	return temp, nil
}

func executeTemplate(temp *template.Template, data interface{}) (string, error) {
	buffer := new(bytes.Buffer)

	err := temp.Execute(buffer, data)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func NewEmail(emailForm models.SendEmailForm) (Email, error) {
	subject := emailForm.Subject
	from := newFrom()
	dest := newDest(emailForm.DestinationName, emailForm.DestinationEmail)
	headers := setHeaders(from.Address, dest.Address, subject)

	html, err := setAndParseTemplate(emailForm.TemplatePath)
	if err != nil {
		return Email{}, err
	}

	emailTemplate, err := executeTemplate(html, emailForm.TemplateData)
	if err != nil {
		return Email{}, err
	}

	var message string

	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += emailTemplate

	return Email{
		Headers: headers,
		Subject: subject,
		From:    from,
		Dest:    dest,
		Message: message,
	}, nil
}
