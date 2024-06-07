package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

func (h Handlers) VerifyMagicLinkHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("staring verify magic link handler")

	slog.Info(fmt.Sprintf("verb=%v", r.Method))

	verifyMagicLinkTemplate := "verify-magic-link"
	templateData := struct {
		IsLoggedIn bool
	}{
		IsLoggedIn: false,
	}

	queryTokenValue := r.URL.Query().Get("token")
	slog.Info(fmt.Sprintf("queryTokenValue=%v", queryTokenValue))

	if queryTokenValue == "" {
		h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
		return
	}

	pendingCookieToken, err := r.Cookie("pendingCookieToken")
	if err != nil {
		h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
		return
	}

	slog.Info(fmt.Sprintf("pendingCookieToken=%v", pendingCookieToken))

	pendingCookieTokenValue := pendingCookieToken.Value
	if pendingCookieTokenValue == "" {
		h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
		return
	}

	slog.Info(fmt.Sprintf("pendingCookieTokenValue=%v", pendingCookieTokenValue))

	pending, err := h.DataService.VerifyPendingUserSession(queryTokenValue, pendingCookieTokenValue)
	if err != nil {
		h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
		return
	}

	slog.Info(fmt.Sprintf("pending=%v", pending))

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

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300, // seconds
	}

	slog.Info(fmt.Sprintf("%v", cookieSessionId))

	http.SetCookie(w, &cookieSessionId)

	cookieToken := http.Cookie{
		Name:     "cookieToken",
		Value:    cookieTokenValue,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300, // seconds
	}

	slog.Info(fmt.Sprintf("%v", cookieToken))
	http.SetCookie(w, &cookieToken)

	templateData.IsLoggedIn = true
	h.templates.ExecuteTemplate(w, verifyMagicLinkTemplate, templateData)
}
