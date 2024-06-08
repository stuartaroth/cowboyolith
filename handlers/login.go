package handlers

import (
	"log/slog"
	"net/http"
)

func (h Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "TRACE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := h.templates.ExecuteTemplate(w, "login", nil)
	if err != nil {
		slog.Error("h.templates.ExecuteTemplate(w, \"login\", templateData)", err)
		return
	}
}
