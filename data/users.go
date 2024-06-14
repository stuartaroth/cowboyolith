package data

import "database/sql"

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
