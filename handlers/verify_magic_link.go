package handlers

import (
	"github.com/google/uuid"
	"net/http"
)

func (h Handlers) VerifyMagicLinkHandler(w http.ResponseWriter, r *http.Request) {
	verifyMagicLinkTemplate := "verify-magic-link"
	templateData := struct {
		IsLoggedIn bool
		IndexUrl   string
	}{
		IsLoggedIn: false,
		IndexUrl:   h.WebServerUrl,
	}

	queryTokenValue := r.URL.Query().Get("token")

	if queryTokenValue == "" {
		h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
		return
	}

	pendingCookieToken, err := r.Cookie("pendingCookieToken")
	if err != nil {
		h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
		return
	}

	pendingCookieTokenValue := pendingCookieToken.Value
	if pendingCookieTokenValue == "" {
		h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
		return
	}

	clearCookie(w, "pendingCookieToken")

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

	setCookie(w, "cookieSessionId", sessionIdValue, 300)
	setCookie(w, "cookieToken", cookieTokenValue, 300)

	templateData.IsLoggedIn = true
	h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
}
