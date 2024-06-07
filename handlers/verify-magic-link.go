package handlers

import (
	"github.com/google/uuid"
	"net/http"
)

func (h Handlers) VerifyMagicLinkHandler(w http.ResponseWriter, r *http.Request) {
	verifyMagicLinkTemplate := "verify-magic-link"
	templateData := struct {
		IsLoggedIn bool
	}{
		IsLoggedIn: false,
	}

	queryTokenValue := r.URL.Query().Get("token")

	if queryTokenValue == "" {
		h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
		return
	}

	pendingCookieToken, err := r.Cookie("pendingCookieToken")
	if err != nil {
		h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
	}

	pendingCookieTokenValue := pendingCookieToken.Value
	if pendingCookieTokenValue == "" {
		h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
		return
	}

	pending, err := h.DataService.VerifyPendingUserSession(queryTokenValue, pendingCookieTokenValue)
	if err != nil {
		h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
		return
	}

	sessionIdValue := uuid.NewString()
	cookieTokenValue := uuid.NewString()

	err = h.DataService.DeletePendingUserSession(pending.Id)
	if err != nil {
		h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
		return
	}

	err = h.DataService.CreateUserSession(pending.UserId, sessionIdValue, cookieTokenValue)
	if err != nil {
		h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
		return
	}

	cookieSessionId := http.Cookie{
		Name:  "cookieSessionId",
		Value: sessionIdValue,
	}
	http.SetCookie(w, &cookieSessionId)

	cookieToken := http.Cookie{
		Name:  "cookieToken",
		Value: cookieTokenValue,
	}
	http.SetCookie(w, &cookieToken)

	templateData.IsLoggedIn = true
	h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
}
