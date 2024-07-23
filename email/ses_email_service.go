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
	webServerUrl   string
	emailClient    *ses.Client
	templates      *template.Template
	shouldSend     bool
	sendingAddress string
}

func NewSesEmailService(webServerUrl string, templates *template.Template, shouldSend bool, sendingAddress string) (SesEmailService, error) {
	var emailClient *ses.Client

	if shouldSend {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return SesEmailService{}, err
		}

		emailClient = ses.NewFromConfig(cfg)
	}

	return SesEmailService{
		webServerUrl:   webServerUrl,
		emailClient:    emailClient,
		templates:      templates,
		shouldSend:     shouldSend,
		sendingAddress: sendingAddress,
	}, nil
}

func (s SesEmailService) SendMagicCode(email, magicCode string) (string, error) {
	emailDestination := types.Destination{
		ToAddresses: []string{email},
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

	sendMailInput := ses.SendEmailInput{
		Destination: &emailDestination,
		Message:     &emailMessage,
		Source:      &s.sendingAddress,
		ReturnPath:  &s.sendingAddress,
	}

	if s.shouldSend {
		_, err = s.emailClient.SendEmail(context.TODO(), &sendMailInput)
		if err != nil {
			return "", err
		}
	}

	return htmlString, nil
}
