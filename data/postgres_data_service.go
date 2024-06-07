package data

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"os"
)

type postgresDataService struct {
	db *sql.DB
}

func NewPostgresDataService() (DataService, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DATABASE")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	sslmode := os.Getenv("POSTGRES_SSL_MODE")

	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", host, port, dbname, user, password, sslmode)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return postgresDataService{
		db: db,
	}, nil
}

func (p postgresDataService) GetAllUsers() ([]User, error) {
	emptyUsers := make([]User, 0)
	users := make([]User, 0)

	rows, err := p.db.Query("select id, email, is_admin, inserted_at from users;")
	if err != nil {
		return emptyUsers, err
	}

	for rows.Next() {
		user, err := ScanUser(rows)
		if err != nil {
			return emptyUsers, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (p postgresDataService) CreatePendingUserSession(userId, id, cookieTokenValue, ipAddress, userAgent string) error {
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

func (p postgresDataService) VerifyPendingUserSession(id, cookieTokenValue string) (PendingUserSession, error) {
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

func (p postgresDataService) CreateUserSession(userId, id, cookieTokenValue string) error {
	salt := uuid.NewString()
	hashedCookieTokenValue, err := hash(cookieTokenValue, salt)
	if err != nil {
		slog.Error("hash", err)
		return err
	}

	_, err = p.db.Exec("insert into user_sessions (id, user_id, hashed_cookie_token, salt) values ($1, $2, $3, $4);", id, userId, hashedCookieTokenValue, salt)
	if err != nil {
		return err
	}

	return nil
}

func (p postgresDataService) GetUserByEmail(email string) (User, error) {
	var u User
	rows, err := p.db.Query("select id, email, is_admin, inserted_at from users where email = $1;", email)
	if err != nil {
		return u, err
	}

	for rows.Next() {

		u, err = ScanUser(rows)
		if err != nil {
			return u, err
		}
	}

	return u, nil
}

func (p postgresDataService) DeletePendingUserSession(id string) error {
	_, err := p.db.Exec("delete from pending_user_sessions where id = $1", id)
	if err != nil {
		slog.Error("DeletePendingUserSession", err)
	}

	return err
}

func ScanUser(rows *sql.Rows) (User, error) {
	var u User
	err := rows.Scan(&u.Id, &u.Email, &u.IsAdmin, &u.InsertedAt)
	if err != nil {
		return u, err
	}

	return u, nil
}

func ScanPendingUserSession(rows *sql.Rows) (PendingUserSession, error) {
	var u PendingUserSession
	err := rows.Scan(&u.Id, &u.UserId, &u.HashedCookieToken, &u.Salt, &u.IpAddress, &u.UserAgent, &u.InsertedAt)
	if err != nil {
		return u, err
	}

	return u, nil
}

func hash(key, salt string) (string, error) {
	bits, err := bcrypt.GenerateFromPassword(getKeySaltBits(key, salt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bits), nil
}

func compareHash(hashedKey, key, salt string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedKey), getKeySaltBits(key, salt))
}

func getKeySaltBits(key, salt string) []byte {
	bits := []byte(fmt.Sprintf("%v%v", key, salt))
	return bits
}
