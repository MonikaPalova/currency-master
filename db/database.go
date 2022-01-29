package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
)

const (
	createTablesFile = "./sql/create_tables.sql"
)

type Database struct {
	conn *sql.DB

	UsersDBHandler *UsersDBHandler
}

func NewDB(user, password string) (*Database, error) {
	dbConnStr := fmt.Sprintf("%s:%s@/?multiStatements=true&parseTime=true", user, password)
	conn, err := sql.Open("mysql", dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to database: %v", err)
	}

	createTablesQuery, err := ioutil.ReadFile(createTablesFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file with create tables queries: %v", err)
	}

	if _, err := conn.Exec(string(createTablesQuery)); err != nil {
		return nil, fmt.Errorf("request to create tables in db failed: %v", err)
	}

	return &Database{conn: conn, UsersDBHandler: &UsersDBHandler{conn: conn}}, nil
}
