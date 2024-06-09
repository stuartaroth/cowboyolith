package handlers

import (
	"github.com/stuartaroth/cowboyolith/constants"
	"github.com/stuartaroth/cowboyolith/data"
	"log/slog"
	"net/http"
)

func (h Handlers) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(constants.User).(data.User)
	sessionId := r.Context().Value(constants.SessionId).(string)

	sessions, err := h.DataService.GetAllUserSessions(user.Id)
	if err != nil {
		sessions = []data.UserSession{}
	}

	templateData := struct {
		Sessions  []data.UserSession
		SessionId string
	}{
		Sessions:  sessions,
		SessionId: sessionId,
	}

	err = h.templates.ExecuteTemplate(w, "sessions", templateData)
	if err != nil {
		slog.Error("h.templates.ExecuteTemplate(w, \"sessions\", templateData)", err)
		return
	}
}
