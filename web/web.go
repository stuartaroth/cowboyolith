package web

import (
	"errors"
	"fmt"
	"github.com/stuartaroth/cowboyolith/data"
	"github.com/stuartaroth/cowboyolith/email"
	"github.com/stuartaroth/cowboyolith/handlers"
	"html/template"
	"net/http"
	"os"
	"strings"
)

func GetWebServerConfig() (string, string, string, string, error) {
	configErrors := []string{}

	webServerUrl := os.Getenv("WEB_SERVER_URL")
	if webServerUrl == "" {
		configErrors = append(configErrors, "WEB_SERVER_URL")
	}

	webServerPort := os.Getenv("WEB_SERVER_PORT")
	if webServerPort == "" {
		configErrors = append(configErrors, "WEB_SERVER_PORT")
	}

	webServerCertFile := os.Getenv("WEB_SERVER_CERT_FILE")
	if webServerCertFile == "" {
		configErrors = append(configErrors, "WEB_SERVER_CERT_FILE")
	}

	webServerKeyFile := os.Getenv("WEB_SERVER_KEY_FILE")
	if webServerKeyFile == "" {
		configErrors = append(configErrors, "WEB_SERVER_KEY_FILE")
	}

	var err error

	if len(configErrors) != 0 {
		missingEnvVariables := strings.Join(configErrors, ",")
		err = errors.New(fmt.Sprintf("following variables not set (%v)", missingEnvVariables))
	}

	return webServerUrl, webServerPort, webServerCertFile, webServerKeyFile, err
}

func GetTemplates(pattern string) (*template.Template, error) {
	return template.ParseGlob(pattern)
}

func GetMux(staticDirectory string, templates *template.Template, dataService data.DataService, emailService email.EmailService) (*http.ServeMux, error) {
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

	// all logged in can use myHandlers.Pre
	myMux.HandleFunc("/", myHandlers.Pre(myHandlers.IndexHandler))
	myMux.HandleFunc("/logout", myHandlers.Pre(myHandlers.LogoutHandler))
	myMux.HandleFunc("/sessions", myHandlers.Pre(myHandlers.SessionsHandler))
	myMux.HandleFunc("POST /sessions/{id}/delete", myHandlers.Pre(myHandlers.DeleteSessionHandler))

	return myMux, nil
}
