package handlers

import (
	"github.com/google/uuid"
	"github.com/stuartaroth/cowboyolith/constants"
	"net/http"
)

func (h Handlers) VerifyMagicLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodTrace {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

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

	pendingCookieToken, err := r.Cookie(constants.PendingCookieToken)
	if err != nil {
		h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
		return
	}

	pendingCookieTokenValue := pendingCookieToken.Value
	if pendingCookieTokenValue == "" {
		h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
		return
	}

	clearCookie(w, constants.PendingCookieToken)

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

	err = h.DataService.CreateUserSession(pending.UserId, sessionIdValue, cookieTokenValue, pending.UserAgent)
	if err != nil {
		h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
		return
	}

	setCookie(w, constants.CookieSessionId, sessionIdValue, constants.OneWeekInSeconds)
	setCookie(w, constants.CookieToken, cookieTokenValue, constants.OneWeekInSeconds)

	templateData.IsLoggedIn = true
	h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
}
