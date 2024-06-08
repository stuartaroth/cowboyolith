package handlers

import (
	"net/http"
)

func (h Handlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {

	sessionId, err := getCookie(r, "cookieSessionId")
	if err == nil {
		h.deleteUserSession(sessionId)
	}

	clearCookie(w, "cookieSessionId")
	clearCookie(w, "cookieToken")
	redirectToLogin(w, r)
}

func (h Handlers) deleteUserSession(sessionId string) {
	go func() {
		h.DataService.DeleteUserSession(sessionId)
	}()
}
