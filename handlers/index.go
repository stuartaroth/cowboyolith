package handlers

import (
	"github.com/stuartaroth/cowboyolith/data"
	"log/slog"
	"net/http"
)

func (h Handlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.templates.ExecuteTemplate(w, "404", nil)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	user := r.Context().Value("user").(data.User)

	templateData := struct {
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
