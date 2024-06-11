package main

import (
	"fmt"
	"github.com/stuartaroth/cowboyolith/config"
	"github.com/stuartaroth/cowboyolith/data"
	"github.com/stuartaroth/cowboyolith/email"
	"github.com/stuartaroth/cowboyolith/web"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	webServerUrl, webServerPort, webServerCertFile, webServerKeyFile, err := config.GetWebServerConfig()
	if err != nil {
		slog.Error("GetWebServerConfig", err)
		os.Exit(1)
	}

	host, port, dbname, user, password, sslmode, err := config.GetDataServiceConfig()
	if err != nil {
		slog.Error("GetDataServiceConfig", err)
		os.Exit(1)
	}

	dataService, err := data.NewPostgresDataService(host, port, dbname, user, password, sslmode)
	if err != nil {
		slog.Error("NewPostgresDataService", err)
		os.Exit(1)
	}

	templates, err := config.GetTemplates("templates/*")
	if err != nil {
		slog.Error("GetTemplates", err)
		os.Exit(1)
	}

	sendEmails := config.GetSendEmails()
	emailService, err := email.NewSesEmailService(webServerUrl, templates, sendEmails)
	if err != nil {
		slog.Error("NewSesEmailService", err)
		os.Exit(1)
	}

	myMux, err := web.GetMux("./static", templates, dataService, emailService)
	if err != nil {
		slog.Error("GetMux", err)
		os.Exit(1)
	}

	slog.Info(fmt.Sprintf("%v see you, space cowboy", webServerUrl))
	err = http.ListenAndServeTLS(fmt.Sprintf(":%v", webServerPort), webServerCertFile, webServerKeyFile, myMux)
	if err != nil {
		slog.Error("http.ListenAndServe", err)
		os.Exit(1)
	}
}
