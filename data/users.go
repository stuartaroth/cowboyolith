package data

import (
	"database/sql"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/stuartaroth/cowboyolith/data/gen/cowboyolith/public/model"
	. "github.com/stuartaroth/cowboyolith/data/gen/cowboyolith/public/table"
)

func JetUserToUser(jetUser model.Users) User {
	return User{
		Id:         jetUser.ID.String(),
		Email:      jetUser.Email,
		IsAdmin:    *jetUser.IsAdmin,
		InsertedAt: jetUser.InsertedAt.String(),
	}
}

func (p PostgresDataService) GetAllUsers() ([]User, error) {
	emptyUsers := make([]User, 0)
	users := make([]User, 0)

	jetUsers := []model.Users{}
	jetStatement := SELECT(Users.ID, Users.Email, Users.IsAdmin, Users.InsertedAt).FROM(Users)

	err := jetStatement.Query(p.db, &jetUsers)
	if err != nil {
		return emptyUsers, err
	}

	for _, jetUser := range jetUsers {
		users = append(users, JetUserToUser(jetUser))
	}

	return users, nil
}

func (p PostgresDataService) GetUserByEmail(email string) (User, error) {
	var u User

	var jetUser model.Users
	jetStatement := SELECT(Users.ID, Users.Email, Users.IsAdmin, Users.InsertedAt).FROM(Users).WHERE(Users.Email.EQ(Text(email)))

	err := jetStatement.Query(p.db, &jetUser)
	if err != nil {
		return u, err
	}

	return JetUserToUser(jetUser), nil
}

func (p PostgresDataService) GetUserById(id string) (User, error) {
	var u User
	parsedUuid, err := uuid.Parse(id)
	if err != nil {
		return u, err
	}

	var jetUser model.Users

	jetStatement := SELECT(Users.ID, Users.Email, Users.IsAdmin, Users.InsertedAt).FROM(Users).WHERE(Users.ID.EQ(UUID(parsedUuid)))

	err = jetStatement.Query(p.db, &jetUser)
	if err != nil {
		return u, err
	}

	return JetUserToUser(jetUser), nil
}

func (p PostgresDataService) CreateUser(dbTx *sql.Tx, id, email string, isAdmin bool) error {
	var err error

	parsedUuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	jetUser := model.Users{
		ID:      parsedUuid,
		Email:   email,
		IsAdmin: &isAdmin,
	}

	jetStatement := Users.INSERT(Users.ID, Users.Email, Users.IsAdmin).MODEL(jetUser)

	if dbTx != nil {
		_, err = jetStatement.Exec(dbTx)
	} else {
		_, err = jetStatement.Exec(p.db)
	}

	return err
}
