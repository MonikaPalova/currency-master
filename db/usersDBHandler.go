package db

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/MonikaPalova/currency-master/httputils"
	"github.com/MonikaPalova/currency-master/model"
)

type UsersDBHandler struct {
	conn *sql.DB
}

func (u *UsersDBHandler) GetAll() ([]model.User, *httputils.HttpError) {
	query := "SELECT username, email, phone_number FROM USERS;"
	rows, err := u.conn.Query(query)
	if err != nil {
		return nil, &httputils.HttpError{Err: err, Message: "could not retrieve users from database", StatusCode: http.StatusInternalServerError}
	}

	return deserializeUsers(rows)
}

func (u *UsersDBHandler) Create(user model.User) (*model.User, *httputils.HttpError) {
	stmt := "INSERT INTO USERS (username, email, phone_number, password) VALUES (?,?,?,?);"
	insertStmt, err := u.conn.Prepare(stmt)
	if err != nil {
		return nil, &httputils.HttpError{Err: err, Message: "error when preparing insert statement for user in database", StatusCode: http.StatusInternalServerError}
	}
	defer insertStmt.Close()

	if _, err = insertStmt.Exec(user.Username, user.Email, user.PhoneNumber, user.Password); err != nil {
		return nil, &httputils.HttpError{Err: err, Message: "error when inserting user in database", StatusCode: http.StatusInternalServerError}
	}
	user.Password = ""
	return &user, nil
}

func (u *UsersDBHandler) GetByUsername(username string) (*model.User, *httputils.HttpError) {
	query := "SELECT username, email, phone_number FROM USERS where username=?;"
	row := u.conn.QueryRow(query, username)

	var user model.User
	if err := row.Scan(&user.Username, &user.Email, &user.PhoneNumber); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &httputils.HttpError{Err: err, Message: "user with username [" + username + "] doesn't exist", StatusCode: http.StatusNotFound}
		}
		return nil, &httputils.HttpError{Err: err, Message: "could not read user row", StatusCode: http.StatusInternalServerError}
	}

	return &user, nil
}

func deserializeUsers(rows *sql.Rows) ([]model.User, *httputils.HttpError) {
	users := []model.User{}
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Username, &user.Email, &user.PhoneNumber); err != nil {
			return nil, &httputils.HttpError{Err: err, Message: "could not read user row", StatusCode: http.StatusInternalServerError}
		}

		users = append(users, user)
	}

	return users, nil
}
