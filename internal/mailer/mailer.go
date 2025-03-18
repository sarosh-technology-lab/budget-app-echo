package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
)

//go:embed templates
var templateFS embed.FS

type Mailer struct {
	dialer *gomail.Dialer
	sender string //"no-reply@budget.com""
	Logger echo.Logger
}

type EmailData struct {
	AppName string
	Subject string
	Meta interface{}
}

func AppMailer(logger echo.Logger) Mailer{
	mailHost := os.Getenv("MAIL_HOST")
	mailPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		logger.Fatal(err)
	}
	mailUser := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")
	mailSender := os.Getenv("MAIL_SENDER")
	
	dialer := gomail.NewDialer(mailHost, mailPort, mailUser, mailPassword)
	return Mailer{
		dialer: dialer,
		sender: mailSender,
		Logger: logger,
	}
}

func (mailer *Mailer) Send(receipient string, templateFile string, data EmailData) error {
	absolutePath := filepath.Join("templates", templateFile)
	tmpl, err := template.ParseFS(templateFS, absolutePath)
	if err != nil {
		mailer.Logger.Error(err)
		return err
	}

	data.AppName = os.Getenv("APP_NAME")

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		mailer.Logger.Error(err)
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		mailer.Logger.Error(err)
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("To", receipient)
	message.SetHeader("From", mailer.sender)
	message.SetHeader("Subject", subject.String())
	message.SetBody("text/html", htmlBody.String())

	err = mailer.dialer.DialAndSend(message)
	if err != nil {
		mailer.Logger.Error(err)
	}

	return nil
}

