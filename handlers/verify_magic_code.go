package handlers

import (
	"github.com/google/uuid"
	"github.com/stuartaroth/cowboyolith/constants"
	"net/http"
)

func (h Handlers) VerifyMagicCodeHandler(w http.ResponseWriter, r *http.Request) {
	magicCode := r.FormValue(constants.MagicCode)

	if magicCode == "" {
		redirectToIndex(w, r)
		return
	}

	pendingCookieToken, err := r.Cookie(constants.PendingCookieToken)
	if err != nil {
		redirectToIndex(w, r)
		return
	}

	pendingCookieTokenValue := pendingCookieToken.Value
	if pendingCookieTokenValue == "" {
		redirectToIndex(w, r)
		return
	}

	clearCookie(w, constants.PendingCookieToken)

	pending, err := h.DataService.VerifyPendingUserSession(magicCode, pendingCookieTokenValue)
	if err != nil {
		redirectToIndex(w, r)
		return
	}

	sessionIdValue := uuid.NewString()
	cookieTokenValue := uuid.NewString()

	dbTx, err := h.DataService.GetDbTransaction()
	if err != nil {
		redirectToIndex(w, r)
		return
	}

	err = h.DataService.DeletePendingUserSession(dbTx, pending.Id)
	if err != nil {
		dbTx.Rollback()
		redirectToIndex(w, r)
		return
	}

	err = h.DataService.CreateUserSession(dbTx, pending.UserId, sessionIdValue, cookieTokenValue, pending.UserAgent)
	if err != nil {
		dbTx.Rollback()
		redirectToIndex(w, r)
		return
	}

	err = dbTx.Commit()
	if err != nil {
		redirectToIndex(w, r)
		return
	}

	setCookie(w, constants.CookieSessionId, sessionIdValue, constants.OneWeekInSeconds)
	setCookie(w, constants.CookieToken, cookieTokenValue, constants.OneWeekInSeconds)
	redirectToIndex(w, r)
}
