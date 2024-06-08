package handlers

import (
	"errors"
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
	cookieSessionId, err := r.Cookie("cookieSessionId")
	if err != nil {
		return "", "", err
	}

	if cookieSessionId.Value == "" {
		return "", "", errors.New("empty cookieSessionId")
	}

	cookieToken, err := r.Cookie("cookieToken")
	if err != nil {
		return "", "", err
	}

	if cookieToken.Value == "" {
		return "", "", errors.New("empty cookieToken")
	}

	return cookieSessionId.Value, cookieToken.Value, nil
}

func (h Handlers) GetCurrentUser(r *http.Request) (data.User, error) {
	sessionId, token, err := GetSessionInfo(r)
	if err != nil {
		return data.User{}, err
	}

	return h.DataService.VerifyUserSession(sessionId, token)
}

func redirectIfLoggedOut(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
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
