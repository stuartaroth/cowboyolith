package handlers

import (
	"net/http"
)

func (h Handlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	clearCookie(w, "cookieSessionId")
	clearCookie(w, "cookieToken")
	redirectIfLoggedOut(w, r)
}
