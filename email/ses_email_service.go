package email

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"html/template"
	"log/slog"
	"os"
)

type sesEmailService struct {
	webServerUrl string
	emailClient  *ses.Client
	templates    *template.Template
}

func NewSesEmailService(logger *slog.Logger, templates *template.Template) (EmailService, error) {
	webServerUrl := os.Getenv("WEB_SERVER_URL")
	if webServerUrl == "" {
		return sesEmailService{}, errors.New("env WEB_SERVER_URL must be set")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return sesEmailService{}, err
	}

	emailClient := ses.NewFromConfig(cfg)

	return sesEmailService{
		webServerUrl: webServerUrl,
		emailClient:  emailClient,
		templates:    templates,
	}, nil
}

func (s sesEmailService) SendMagicLink(email, queryToken string) error {
	emailDestination := types.Destination{
		BccAddresses: []string{},
		CcAddresses:  []string{},
		ToAddresses:  []string{email},
	}

	charset := "UTF-8"

	verifyMagicLinkUrl := fmt.Sprintf("%v/verify-magic-link?token=%v", s.webServerUrl, queryToken)

	templateData := struct {
		VerifyMagicLinkUrl string
	}{
		VerifyMagicLinkUrl: verifyMagicLinkUrl,
	}

	bits := new(bytes.Buffer)
	err := s.templates.ExecuteTemplate(bits, "magic-link-email", templateData)
	if err != nil {
		return err
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

	subjectLine := "Cowboyolith Magic Link"

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

	_, err = s.emailClient.SendEmail(context.TODO(), &sendMailInput)
	if err != nil {
		return err
	}

	return nil
}
