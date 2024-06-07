package main

import (
	"github.com/stuartaroth/cowboyolith/data"
	"github.com/stuartaroth/cowboyolith/email"
	"github.com/stuartaroth/cowboyolith/handlers"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	templates, err := template.ParseGlob("templates/*")
	if err != nil {
		slogger.Error("templates", err)
		os.Exit(1)
	}

	dataService, err := data.NewPostgresDataService(slogger)
	if err != nil {
		slogger.Error("data.NewPostgresDataService", err)
		os.Exit(1)
	}

	emailService, err := email.NewSesEmailService(slogger, templates)
	if err != nil {
		slogger.Error("email.NewSesEmailService", err)
		os.Exit(1)
	}

	myHandlers, err := handlers.NewHandlers(slogger, dataService, emailService, templates)
	if err != nil {
		slogger.Error("handlers.NewHandlers", err)
		os.Exit(1)
	}

	http.HandleFunc("/", myHandlers.IndexHandler)
	http.HandleFunc("POST /email-request", myHandlers.EmailRequestHandler)
	http.HandleFunc("/verify-magic-link", myHandlers.VerifyMagicLinkHandler)

	slogger.Info("see you, space cowboy")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		slogger.Error("http.ListenAndServe", err)
		os.Exit(1)
	}
}
