package tests

import (
	"github.com/stuartaroth/cowboyolith/config"
	"github.com/stuartaroth/cowboyolith/email"
	"strings"
	"testing"
)

func TestSendMagicLink(t *testing.T) {
	templates, err := config.GetTemplates("../templates/*")
	if err != nil {
		t.Fatal("error GetTemplates()", err)
	}

	url := "url"
	magicCode := "magicCode"

	emailService, err := email.NewSesEmailService(url, templates, false, "sending@gmail.com")
	if err != nil {
		t.Fatal("error NewSesEmailService", err)
	}

	html, err := emailService.SendMagicCode("", magicCode)
	if err != nil {
		t.Fatal("error SendMagicCode", err)
	}

	stringsToCheck := []string{
		"Your code is valid for five minutes and must be used in the requesting browser",
		magicCode,
	}

	for _, i := range stringsToCheck {
		if !strings.Contains(html, i) {
			t.Fatal("html does not contain expected string:", i)
		}
	}
}
