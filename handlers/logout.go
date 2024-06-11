package handlers

import (
	"github.com/stuartaroth/cowboyolith/constants"
	"github.com/stuartaroth/cowboyolith/data"
	"log/slog"
	"net/http"
)

func (h Handlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(constants.User).(data.User)
	sessionId := r.Context().Value(constants.SessionId).(string)
	err := h.DataService.DeleteUserSession(user.Id, sessionId)
	if err != nil {
		slog.Error("error in logout handler", err)
	}

	redirectToLogin(w, r)
}
