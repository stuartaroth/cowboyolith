package handlers

import (
	"github.com/google/uuid"
	"github.com/stuartaroth/cowboyolith/constants"
	"log/slog"
	"net/http"
)

func (h Handlers) EmailRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodTrace {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userEmail := r.FormValue(constants.Email)
	cookieTokenValue := uuid.NewString()

	setCookie(w, constants.PendingCookieToken, cookieTokenValue, 120)

	h.createPendingUserSessionAndEmail(r, userEmail, cookieTokenValue)

	h.templates.ExecuteTemplate(w, "email-request", nil)
}

func (h Handlers) createPendingUserSessionAndEmail(r *http.Request, userEmail, cookieTokenValue string) {
	go func() {
		user, err := h.DataService.GetUserByEmail(userEmail)
		if err != nil {
			slog.Error("h.DataService.GetUserByEmail", err)
			return
		}

		id := uuid.NewString()
		ipAddress := r.RemoteAddr
		userAgent := r.UserAgent()

		err = h.DataService.CreatePendingUserSession(user.Id, id, cookieTokenValue, ipAddress, userAgent)
		if err != nil {
			slog.Error("h.DataService.CreatePendingUserSession", err)
			return
		}

		_, err = h.EmailService.SendMagicLink(user.Email, id)
		if err != nil {
			slog.Error("h.EmailService.SendMagicLink", err)
			return
		}
	}()
}
