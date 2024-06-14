package config

import (
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"strings"
)

func GetTemplates(pattern string) (*template.Template, error) {
	return template.ParseGlob(pattern)
}

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

func GetDataServiceConfig() (string, string, string, string, string, string, error) {
	configErrors := []string{}

	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		configErrors = append(configErrors, "POSTGRES_HOST")
	}

	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		configErrors = append(configErrors, "POSTGRES_PORT")
	}

	dbname := os.Getenv("POSTGRES_DATABASE")
	if dbname == "" {
		configErrors = append(configErrors, "POSTGRES_DATABASE")
	}

	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		configErrors = append(configErrors, "POSTGRES_USER")
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		configErrors = append(configErrors, "POSTGRES_PASSWORD")
	}

	sslmode := os.Getenv("POSTGRES_SSL_MODE")
	if sslmode == "" {
		configErrors = append(configErrors, "POSTGRES_SSL_MODE")
	}

	return host, port, dbname, user, password, sslmode, nil
}

func GetSendEmails() bool {
	return os.Getenv("SEND_EMAILS") == "true"
}

func GetLogLevel() slog.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}
