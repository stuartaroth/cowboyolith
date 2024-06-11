package data

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type PostgresDataService struct {
	db *sql.DB
}

func GetDataServiceConfig() (string, string, string, string, string, string, error) {
	configErrors := []string{}

	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		configErrors = append(configErrors, "POSTGRES_HOST")
	}

	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		configErrors = append(configErrors, "POSTGRES_PORT")
	}

	dbname := os.Getenv("POSTGRES_DATABASE")
	if dbname == "" {
		configErrors = append(configErrors, "POSTGRES_DATABASE")
	}

	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		configErrors = append(configErrors, "POSTGRES_USER")
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		configErrors = append(configErrors, "POSTGRES_PASSWORD")
	}

	sslmode := os.Getenv("POSTGRES_SSL_MODE")
	if sslmode == "" {
		configErrors = append(configErrors, "POSTGRES_SSL_MODE")
	}

	return host, port, dbname, user, password, sslmode, nil
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
