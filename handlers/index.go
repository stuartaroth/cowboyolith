package handlers

import (
	"github.com/stuartaroth/cowboyolith/constants"
	"github.com/stuartaroth/cowboyolith/data"
	"log/slog"
	"net/http"
)

func (h Handlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.templates.ExecuteTemplate(w, "404", map[string]string{"Title": "Not Found"})
		return
	}

	user := r.Context().Value(constants.User).(data.User)

	templateData := struct {
		Title string
		Email string
	}{
		Email: user.Email,
	}

	err := h.templates.ExecuteTemplate(w, "index", templateData)
	if err != nil {
		slog.Error("h.templates.ExecuteTemplate(w, \"index\", templateData)", err)
		return
	}
}
