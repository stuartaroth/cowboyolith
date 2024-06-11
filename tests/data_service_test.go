package tests

import (
	"database/sql"
	"fmt"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/google/uuid"
	"github.com/stuartaroth/cowboyolith/data"
	"os"
	"testing"
)

func SetupDatabaseAndRunCreate() (*embeddedpostgres.EmbeddedPostgres, error) {
	bits, err := os.ReadFile("../data/setup.sql")
	if err != nil {
		return nil, err
	}

	createSql := string(bits)

	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().Port(5433))

	host := "localhost"
	port := "5433"
	dbname := "postgres"
	user := "postgres"
	password := "postgres"
	sslmode := "disable"

	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", host, port, dbname, user, password, sslmode)

	err = postgres.Start()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		postgres.Stop()
		return nil, err
	}

	defer db.Close()

	_, err = db.Exec(createSql)
	if err != nil {
		postgres.Stop()
		return nil, err
	}

	return postgres, nil
}

func TestSetup(t *testing.T) {
	postgres, err := SetupDatabaseAndRunCreate()
	if err != nil {
		t.Fatal(err)
	}

	defer postgres.Stop()

	host := "localhost"
	port := "5433"
	dbname := "postgres"
	user := "postgres"
	password := "postgres"
	sslmode := "disable"

	dataService, err := data.NewPostgresDataService(host, port, dbname, user, password, sslmode)
	if err != nil {
		t.Fatal(err)
	}

	userId := uuid.NewString()
	userEmail := "test@gmail.com"

	err = dataService.CreateUser(userId, userEmail, true)
	if err != nil {
		t.Fatal(err)
	}

	users, err := dataService.GetAllUsers()
	if err != nil {
		t.Fatal("GetAllUsers", err)
	}

	if len(users) != 1 {
		t.Fatal("len(users")
	}

	storedUser1, err := dataService.GetUserByEmail(userEmail)
	if err != nil {
		t.Fatal("GetUserByEmail", err)
	}

	if storedUser1.Id != userId {
		t.Fatal("storedUser1.Id")
	}

	if storedUser1.Email != userEmail {
		t.Fatal("storedUser1.Email")
	}

	if !storedUser1.IsAdmin {
		t.Fatal("storedUser1.IsAdmin")
	}

	storedUser2, err := dataService.GetUserById(userId)
	if err != nil {
		t.Fatal("GetUserById", err)
	}

	if storedUser2.Id != userId {
		t.Fatal("storedUser1.Id")
	}

	if storedUser2.Email != userEmail {
		t.Fatal("storedUser1.Email")
	}

	if !storedUser2.IsAdmin {
		t.Fatal("storedUser1.IsAdmin")
	}
}
