package db

import (
	"database/sql"
	"fmt"

	"github.com/MonikaPalova/currency-master/model"
	"github.com/go-sql-driver/mysql"
)

const (
	SELECT_USER_AND_ASSETS             = "SELECT USERS.username, USERS.email, USER_ASSETS.asset_id, USER_ASSETS.name, USER_ASSETS.quantity FROM USERS LEFT JOIN USER_ASSETS ON USERS.username=USER_ASSETS.username;"
	SELECT_USER_AND_ASSETS_BY_USERNAME = "SELECT USERS.username, USERS.email, USER_ASSETS.asset_id, USER_ASSETS.name, USER_ASSETS.quantity FROM USERS LEFT JOIN USER_ASSETS ON USERS.username=USER_ASSETS.username where USERS.username=?;"
	INSERT_USER                        = "INSERT INTO USERS (username, email, password) VALUES (?,?,?);"
)

type UsersDBHandler struct {
	conn *sql.DB
}

type userAsset struct {
	user     model.User
	assetId  sql.NullString
	name     sql.NullString
	quantity sql.NullFloat64
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
	user.Assets = []model.UserAsset{}
	return &user, nil
}

func (u UsersDBHandler) GetAll() ([]model.User, error) {
	rows, err := u.conn.Query(SELECT_USER_AND_ASSETS)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve users from database, %v", err)
	}

	return deserializeUsers(rows)
}

func (u UsersDBHandler) GetByUsername(username string) (*model.User, error) {
	rows, err := u.conn.Query(SELECT_USER_AND_ASSETS_BY_USERNAME, username)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve user from database, %v", err)
	}

	users, err := deserializeUsers(rows)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return &users[0], nil
}

func deserializeUsers(rows *sql.Rows) ([]model.User, error) {
	usersByUsername := make(map[string]model.User)

	for rows.Next() {
		var asset userAsset
		if err := rows.Scan(&asset.user.Username, &asset.user.Email, &asset.assetId, &asset.name, &asset.quantity); err != nil {
			return nil, fmt.Errorf("could not read user row, %v", err)
		}
		if user, exists := usersByUsername[asset.user.Username]; exists {
			user.Assets = append(user.Assets, model.UserAsset{AssetId: asset.assetId.String, Name: asset.name.String, Quantity: asset.quantity.Float64})
			usersByUsername[asset.user.Username] = user
		} else {
			user := asset.user
			if asset.assetId.Valid {
				user.Assets = []model.UserAsset{{AssetId: asset.assetId.String, Name: asset.name.String, Quantity: asset.quantity.Float64}}
			} else {
				user.Assets = []model.UserAsset{}
			}
			usersByUsername[asset.user.Username] = user
		}
	}

	var users []model.User
	for _, user := range usersByUsername {
		users = append(users, user)
	}

	return users, nil
}
