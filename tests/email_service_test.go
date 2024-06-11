package tests

import (
	"fmt"
	"github.com/stuartaroth/cowboyolith/email"
	"github.com/stuartaroth/cowboyolith/web"
	"strings"
	"testing"
)

func TestSendMagicLink(t *testing.T) {
	templates, err := web.GetTemplates("../templates/*")
	if err != nil {
		t.Fatal("error GetTemplates()", err)
	}

	url := "url"
	queryToken := "queryToken"

	emailService, err := email.NewSesEmailService(url, templates, false)
	if err != nil {
		t.Fatal("error NewSesEmailService", err)
	}

	html, err := emailService.SendMagicLink("", queryToken)
	if err != nil {
		t.Fatal("error SendMagicLink", err)
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
