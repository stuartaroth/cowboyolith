package web

import (
	"github.com/stuartaroth/cowboyolith/data"
	"github.com/stuartaroth/cowboyolith/email"
	"github.com/stuartaroth/cowboyolith/handlers"
	"html/template"
	"net/http"
)

func GetMux(staticDirectory string, templates *template.Template, dataService data.PostgresDataService, emailService email.SesEmailService) (*http.ServeMux, error) {
	myHandlers, err := handlers.NewHandlers(dataService, emailService, templates)
	if err != nil {
		return nil, err
	}

	myMux := http.NewServeMux()
	fs := http.FileServer(http.Dir(staticDirectory))
	myMux.Handle("/static/", http.StripPrefix("/static/", fs))

	// login functions
	myMux.HandleFunc("/login", myHandlers.LoginHandler)

	myMux.HandleFunc("POST /email-request", myHandlers.EmailRequestHandler)
	myMux.HandleFunc("/verify-magic-link", myHandlers.VerifyMagicLinkHandler)

	// all logged in can use myHandlers.Authorized
	myMux.HandleFunc("/", myHandlers.Authorized(myHandlers.IndexHandler))
	myMux.HandleFunc("/logout", myHandlers.Authorized(myHandlers.LogoutHandler))
	myMux.HandleFunc("/sessions", myHandlers.Authorized(myHandlers.SessionsHandler))
	myMux.HandleFunc("POST /sessions/{id}/delete", myHandlers.Authorized(myHandlers.DeleteSessionHandler))

	return myMux, nil
}
