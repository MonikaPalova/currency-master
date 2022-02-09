package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/MonikaPalova/currency-master/model"
	"github.com/go-sql-driver/mysql"
)

const (
	SELECT_USERS            = "SELECT username, email FROM USERS;"
	INSERT_USER             = "INSERT INTO USERS (username, email, password) VALUES (?,?,?);"
	SELECT_USER_BY_USERNAME = "SELECT username, email FROM USERS where username=?;"
)

type UsersDBHandler struct {
	conn *sql.DB
}

func (u UsersDBHandler) GetAll() ([]model.User, error) {
	rows, err := u.conn.Query(SELECT_USERS)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve users from database, %v", err)
	}

	return deserializeUsers(rows)
}

func (u UsersDBHandler) Create(user model.User) (*model.User, error) {
	insertStmt, err := u.conn.Prepare(INSERT_USER)
	if err != nil {
		return nil, fmt.Errorf("error when preparing insert statement for user in database, %v", err)
	}
	defer insertStmt.Close()

	if _, err = insertStmt.Exec(user.Username, user.Email, user.Password); err != nil {
		if err.(*mysql.MySQLError).Number == 1062 {
			return nil, nil
		}
		return nil, fmt.Errorf("error when inserting user in database, %v", err)
	}
	user.Password = ""
	return &user, nil
}

func (u UsersDBHandler) GetByUsername(username string) (*model.User, error) {
	row := u.conn.QueryRow(SELECT_USER_BY_USERNAME, username)

	var user model.User
	if err := row.Scan(&user.Username, &user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not read user row, %v", err)
	}

	return &user, nil
}

func deserializeUsers(rows *sql.Rows) ([]model.User, error) {
	users := []model.User{}
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Username, &user.Email); err != nil {
			return nil, fmt.Errorf("could not read user row, %v", err)
		}

		users = append(users, user)
	}

	return users, nil
}
