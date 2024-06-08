package handlers

import (
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

func (h Handlers) EmailRequestHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := r.FormValue("email")
	cookieTokenValue := uuid.NewString()

	setCookie(w, "pendingCookieToken", cookieTokenValue, 120)

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

		err = h.EmailService.SendMagicLink(user.Email, id)
		if err != nil {
			slog.Error("h.EmailService.SendMagicLink", err)
			return
		}
	}()
}
