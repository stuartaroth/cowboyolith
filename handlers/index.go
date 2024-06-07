package handlers

import (
	"github.com/stuartaroth/cowboyolith/data"
	"log/slog"
	"net/http"
)

func (h Handlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.DataService.GetAllUsers()
	if err != nil {
		users = []data.User{}
	}

	templateData := struct {
		Users      []data.User
		IsLoggedIn bool
	}{
		Users:      users,
		IsLoggedIn: false,
	}

	err = h.templates.ExecuteTemplate(w, "index", templateData)
	if err != nil {
		slog.Error("h.templates.ExecuteTemplate(w, \"index\", templateData)", err)
		return
	}
}
