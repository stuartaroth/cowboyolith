package data

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"log/slog"
)

func ScanUserSession(rows *sql.Rows) (UserSession, error) {
	var u UserSession
	err := rows.Scan(&u.Id, &u.UserId, &u.HashedCookieToken, &u.Salt, &u.UserAgent, &u.InsertedAt)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (p PostgresDataService) CreateUserSession(dbTx *sql.Tx, userId, id, cookieTokenValue, userAgent string) error {
	salt := uuid.NewString()
	hashedCookieTokenValue, err := hash(cookieTokenValue, salt)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	createUserSessionQuery := "insert into user_sessions (id, user_id, hashed_cookie_token, salt, user_agent) values ($1, $2, $3, $4, $5);"

	if dbTx != nil {
		_, err = dbTx.Exec(createUserSessionQuery, id, userId, hashedCookieTokenValue, salt, userAgent)
	} else {
		_, err = p.db.Exec(createUserSessionQuery, id, userId, hashedCookieTokenValue, salt, userAgent)
	}

	return err
}

func (p PostgresDataService) DeleteUserSession(userId, sessionId string) error {
	_, err := p.db.Exec("delete from user_sessions where id = $1 and user_id = $2", sessionId, userId)
	if err != nil {
		slog.Error(err.Error())
	}

	return err
}

func (p PostgresDataService) GetAllUserSessions(userId string) ([]UserSession, error) {
	sessions := make([]UserSession, 0)

	query :=
		`select
		    id, user_id, hashed_cookie_token, salt, user_agent, inserted_at
		from user_sessions
			where user_id = $1;`

	rows, err := p.db.Query(query, userId)
	if err != nil {
		return sessions, err
	}

	for rows.Next() {
		userSession, err := ScanUserSession(rows)
		if err != nil {
			return sessions, err
		}

		sessions = append(sessions, userSession)
	}

	return sessions, nil
}

func (p PostgresDataService) VerifyUserSession(id, token string) (User, error) {
	var user User
	var userSession UserSession
	query :=
		`select 
		    id, user_id, hashed_cookie_token, salt, user_agent, inserted_at
		from user_sessions 
			where id = $1
		limit 1;`

	rows, err := p.db.Query(query, id)
	if err != nil {
		return user, err
	}

	for rows.Next() {
		userSession, err = ScanUserSession(rows)
		if err != nil {
			return user, err
		}
	}

	if userSession.Id == "" {
		return user, errors.New("no result")
	}

	err = compareHash(userSession.HashedCookieToken, token, userSession.Salt)
	if err != nil {
		return user, err
	}

	return p.GetUserById(userSession.UserId)
}
