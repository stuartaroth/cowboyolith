package tests

import (
	"fmt"
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
	queryToken := "queryToken"

	emailService, err := email.NewSesEmailService(url, templates, false)
	if err != nil {
		t.Fatal("error NewSesEmailService", err)
	}

	html, err := emailService.SendMagicCode("", queryToken)
	if err != nil {
		t.Fatal("error SendMagicCode", err)
	}

	stringsToCheck := []string{
		"Click here to login",
		fmt.Sprintf("href='%v/verify-magic-link?token=%v'", url, queryToken),
	}

	for _, i := range stringsToCheck {
		if !strings.Contains(html, i) {
			t.Fatal("html does not contain expected string:", i)
		}
	}
}
