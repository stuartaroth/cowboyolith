package email

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"html/template"
)

type SesEmailService struct {
	webServerUrl string
	emailClient  *ses.Client
	templates    *template.Template
	sendEmails   bool
}

func NewSesEmailService(webServerUrl string, templates *template.Template, sendEmails bool) (SesEmailService, error) {
	var emailClient *ses.Client

	if sendEmails {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return SesEmailService{}, err
		}

		emailClient = ses.NewFromConfig(cfg)
	}

	return SesEmailService{
		webServerUrl: webServerUrl,
		emailClient:  emailClient,
		templates:    templates,
		sendEmails:   sendEmails,
	}, nil
}

func (s SesEmailService) SendMagicCode(email, magicCode string) (string, error) {
	emailDestination := types.Destination{
		BccAddresses: []string{},
		CcAddresses:  []string{},
		ToAddresses:  []string{email},
	}

	charset := "UTF-8"

	templateData := struct {
		MagicCode string
	}{
		MagicCode: magicCode,
	}

	bits := new(bytes.Buffer)
	err := s.templates.ExecuteTemplate(bits, "magic-code-email", templateData)
	if err != nil {
		return "", err
	}

	htmlString := bits.String()

	emailBodyContent := types.Content{
		Data:    &htmlString,
		Charset: &charset,
	}

	emailBody := types.Body{
		Html: &emailBodyContent,
		Text: nil,
	}

	subjectLine := "Cowboyolith Magic Code"

	emailSubject := types.Content{
		Data:    &subjectLine,
		Charset: &charset,
	}

	emailMessage := types.Message{
		Body:    &emailBody,
		Subject: &emailSubject,
	}

	sendingEmailAddress := "stuartaroth@gmail.com"

	sendMailInput := ses.SendEmailInput{
		Destination: &emailDestination,
		Message:     &emailMessage,
		Source:      &sendingEmailAddress,
		ReturnPath:  &sendingEmailAddress,
	}

	if s.sendEmails {
		_, err = s.emailClient.SendEmail(context.TODO(), &sendMailInput)
		if err != nil {
			return "", err
		}
	}

	return htmlString, nil
}
