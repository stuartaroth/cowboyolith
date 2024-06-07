package handlers

import (
	"github.com/stuartaroth/cowboyolith/data"
	"github.com/stuartaroth/cowboyolith/email"
	"html/template"
)

type Handlers struct {
	DataService  data.DataService
	EmailService email.EmailService
	templates    *template.Template
}

func NewHandlers(dataService data.DataService, emailService email.EmailService, templates *template.Template) (Handlers, error) {
	return Handlers{
		DataService:  dataService,
		EmailService: emailService,
		templates:    templates,
	}, nil
}
