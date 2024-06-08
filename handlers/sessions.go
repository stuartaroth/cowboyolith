package handlers

import (
	"github.com/stuartaroth/cowboyolith/data"
	"log/slog"
	"net/http"
)

func (h Handlers) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(data.User)

	sessions, err := h.DataService.GetAllUserSessions(user.Id)
	if err != nil {
		sessions = []data.UserSession{}
	}

	templateData := struct {
		Sessions []data.UserSession
	}{
		Sessions: sessions,
	}

	err = h.templates.ExecuteTemplate(w, "sessions", templateData)
	if err != nil {
		slog.Error("h.templates.ExecuteTemplate(w, \"sessions\", templateData)", err)
		return
	}
}
