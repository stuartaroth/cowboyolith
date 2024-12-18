package data

import (
	"database/sql"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/stuartaroth/cowboyolith/data/gen/cowboyolith/public/model"
	. "github.com/stuartaroth/cowboyolith/data/gen/cowboyolith/public/table"
)

func ScanUser(rows *sql.Rows) (User, error) {
	var u User
	err := rows.Scan(&u.Id, &u.Email, &u.IsAdmin, &u.InsertedAt)
	if err != nil {
		return u, err
	}

	return u, nil
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
		users = append(users, User{
			Id:         jetUser.ID.String(),
			Email:      jetUser.Email,
			IsAdmin:    *jetUser.IsAdmin,
			InsertedAt: jetUser.InsertedAt.String(),
		})
	}

	return users, nil
}

func (p PostgresDataService) GetUserByEmail(email string) (User, error) {
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

func (p PostgresDataService) GetUserById(id string) (User, error) {
	var u User
	rows, err := p.db.Query("select id, email, is_admin, inserted_at from users where id = $1;", id)
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

func (p PostgresDataService) CreateUser(dbTx *sql.Tx, id, email string, isAdmin bool) error {
	var err error
	insertUserQuery := "insert into users (id, email, is_admin) values ($1, $2, $3);"

	if dbTx != nil {
		_, err = dbTx.Exec(insertUserQuery, id, email, isAdmin)
	} else {
		_, err = p.db.Exec(insertUserQuery, id, email, isAdmin)
	}

	return err
}
