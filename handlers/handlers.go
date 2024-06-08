package handlers

import (
	"errors"
	"fmt"
	"github.com/stuartaroth/cowboyolith/data"
	"github.com/stuartaroth/cowboyolith/email"
	"html/template"
	"net/http"
)

type Handlers struct {
	WebServerUrl string
	DataService  data.DataService
	EmailService email.EmailService
	templates    *template.Template
}

func NewHandlers(webServerUrl string, dataService data.DataService, emailService email.EmailService, templates *template.Template) (Handlers, error) {
	return Handlers{
		WebServerUrl: webServerUrl,
		DataService:  dataService,
		EmailService: emailService,
		templates:    templates,
	}, nil
}

func GetSessionInfo(r *http.Request) (string, string, error) {
	cookieSessionIdValue, err := getCookie(r, "cookieSessionId")
	if err != nil {
		return "", "", err
	}

	cookieTokenValue, err := getCookie(r, "cookieToken")
	if err != nil {
		return "", "", err
	}

	return cookieSessionIdValue, cookieTokenValue, nil
}

func (h Handlers) GetCurrentUser(r *http.Request) (data.User, error) {
	sessionId, token, err := GetSessionInfo(r)
	if err != nil {
		return data.User{}, err
	}

	return h.DataService.VerifyUserSession(sessionId, token)
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func getCookie(r *http.Request, name string) (string, error) {
	storedCookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	if storedCookie.Value == "" {
		return "", errors.New(fmt.Sprintf("empty cookie: %v", name))
	}

	return storedCookie.Value, nil
}

func setCookie(w http.ResponseWriter, name, value string, maxAge int) {
	newCookie := http.Cookie{
		Name:     name,
		Value:    value,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	}

	http.SetCookie(w, &newCookie)
}

func clearCookie(w http.ResponseWriter, name string) {
	newCookie := http.Cookie{
		Name:     name,
		Value:    "",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}

	http.SetCookie(w, &newCookie)
}
