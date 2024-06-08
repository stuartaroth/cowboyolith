package main

import (
	"errors"
	"fmt"
	"github.com/stuartaroth/cowboyolith/data"
	"github.com/stuartaroth/cowboyolith/email"
	"github.com/stuartaroth/cowboyolith/handlers"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

func getWebServerConfig() (string, string, string, string, error) {
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

func main() {
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(slogger)

	webServerUrl, webServerPort, webServerCertFile, webServerKeyFile, err := getWebServerConfig()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	templates, err := template.ParseGlob("templates/*")
	if err != nil {
		slog.Error("templates", err)
		os.Exit(1)
	}

	dataService, err := data.NewPostgresDataService()
	if err != nil {
		slog.Error("data.NewPostgresDataService", err)
		os.Exit(1)
	}

	emailService, err := email.NewSesEmailService(webServerUrl, templates)
	if err != nil {
		slog.Error("email.NewSesEmailService", err)
		os.Exit(1)
	}

	myHandlers, err := handlers.NewHandlers(webServerUrl, dataService, emailService, templates)
	if err != nil {
		slog.Error("handlers.NewHandlers", err)
		os.Exit(1)
	}

	// login and logout functions
	http.HandleFunc("/login", myHandlers.LoginHandler)
	http.HandleFunc("/logout", myHandlers.LogoutHandler)
	http.HandleFunc("POST /email-request", myHandlers.EmailRequestHandler)
	http.HandleFunc("/verify-magic-link", myHandlers.VerifyMagicLinkHandler)

	// all logged in can use myHandlers.Pre
	http.HandleFunc("/", myHandlers.Pre(myHandlers.IndexHandler))
	http.HandleFunc("/sessions", myHandlers.Pre(myHandlers.SessionsHandler))

	slog.Info(fmt.Sprintf("%v see you, space cowboy", webServerUrl))
	err = http.ListenAndServeTLS(fmt.Sprintf(":%v", webServerPort), webServerCertFile, webServerKeyFile, nil)
	if err != nil {
		slog.Error("http.ListenAndServe", err)
		os.Exit(1)
	}
}
