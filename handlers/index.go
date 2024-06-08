package handlers

import (
	"log/slog"
	"net/http"
)

func (h Handlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
	user, err := h.GetCurrentUser(r)
	if err != nil {
		redirectIfLoggedOut(w, r)
		return
	}

	templateData := struct {
		Email string
	}{
		Email: user.Email,
	}

	err = h.templates.ExecuteTemplate(w, "index", templateData)
	if err != nil {
		slog.Error("h.templates.ExecuteTemplate(w, \"index\", templateData)", err)
		return
	}
}
