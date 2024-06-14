package data

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type PostgresDataService struct {
	db *sql.DB
}

func NewPostgresDataService(host, port, dbname, user, password, sslmode string) (PostgresDataService, error) {
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", host, port, dbname, user, password, sslmode)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return PostgresDataService{}, err
	}

	return PostgresDataService{
		db: db,
	}, nil
}

func (p PostgresDataService) GetDbTransaction() (*sql.Tx, error) {
	return p.db.Begin()
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
