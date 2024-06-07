package handlers

import (
	"github.com/stuartaroth/cowboyolith/data"
	"github.com/stuartaroth/cowboyolith/email"
	"html/template"
	"log/slog"
)

type Handlers struct {
	slogger      *slog.Logger
	DataService  data.DataService
	EmailService email.EmailService
	templates    *template.Template
}

func NewHandlers(slogger *slog.Logger, dataService data.DataService, emailService email.EmailService, templates *template.Template) (Handlers, error) {
	return Handlers{
		slogger:      slogger,
		DataService:  dataService,
		EmailService: emailService,
		templates:    templates,
	}, nil
}
