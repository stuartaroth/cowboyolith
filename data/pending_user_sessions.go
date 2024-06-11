package data

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"log/slog"
)

func ScanPendingUserSession(rows *sql.Rows) (PendingUserSession, error) {
	var u PendingUserSession
	err := rows.Scan(&u.Id, &u.UserId, &u.HashedCookieToken, &u.Salt, &u.IpAddress, &u.UserAgent, &u.InsertedAt)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (p PostgresDataService) CreatePendingUserSession(userId, id, cookieTokenValue, ipAddress, userAgent string) error {
	salt := uuid.NewString()
	hashedCookieTokenValue, err := hash(cookieTokenValue, salt)
	if err != nil {
		return err
	}

	_, err = p.db.Exec("insert into pending_user_sessions (id, user_id, hashed_cookie_token, salt, ip_address, user_agent) values ($1, $2, $3, $4, $5, $6);", id, userId, hashedCookieTokenValue, salt, ipAddress, userAgent)
	if err != nil {
		return err
	}

	return nil
}

func (p PostgresDataService) VerifyPendingUserSession(id, cookieTokenValue string) (PendingUserSession, error) {
	var u PendingUserSession
	query :=
		`select 
		    id, user_id, hashed_cookie_token, salt, ip_address, user_agent, inserted_at 
		from pending_user_sessions 
			where id = $1
			and inserted_at > (timezone('utc', now()) - interval '5 minute')
		limit 1;`

	rows, err := p.db.Query(query, id)
	if err != nil {
		return u, err
	}

	for rows.Next() {
		u, err = ScanPendingUserSession(rows)
		if err != nil {
			return u, err
		}
	}

	if u.Id == "" {
		return u, errors.New("no result")
	}

	err = compareHash(u.HashedCookieToken, cookieTokenValue, u.Salt)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (p PostgresDataService) DeletePendingUserSession(id string) error {
	_, err := p.db.Exec("delete from pending_user_sessions where id = $1", id)
	if err != nil {
		slog.Error("DeletePendingUserSession", err)
	}

	return err
}
