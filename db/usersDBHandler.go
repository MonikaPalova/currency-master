package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/MonikaPalova/currency-master/model"
	"github.com/go-sql-driver/mysql"
)

const (
	selectUserAndAssets           = "SELECT USERS.username, USERS.email, USERS.usd, USER_ASSETS.asset_id, USER_ASSETS.name, USER_ASSETS.quantity FROM USERS LEFT JOIN USER_ASSETS ON USERS.username=USER_ASSETS.username;"
	selectUserAndAssetsByUsername = "SELECT USERS.username, USERS.email, USERS.usd, USER_ASSETS.asset_id, USER_ASSETS.name, USER_ASSETS.quantity FROM USERS LEFT JOIN USER_ASSETS ON USERS.username=USER_ASSETS.username where USERS.username=?;"
	insertUser                    = "INSERT INTO USERS (username, email, password,usd) VALUES (?,?,?,?);"
	selectUser                    = "SELECT username, email, usd FROM USERS where username=?;"
	updateUserUSD                 = "UPDATE USERS SET usd = ? WHERE username=?;"
	existsUser                    = "SELECT COUNT(1) FROM USERS WHERE username=? AND password=?;"
)

// Handles sql operations to USERS table.
type UsersDBHandler struct {
	conn *sql.DB
}

type userAsset struct {
	user     model.User
	assetId  sql.NullString
	name     sql.NullString
	quantity sql.NullFloat64
}

// Saves new user in database.
// Returns same user with blank password and empty assets array, if successful.
// Returns nil user if user already exists.
// Returns error on database query error
func (u UsersDBHandler) Create(user model.User) (*model.User, error) {
	insertStmt, err := u.conn.Prepare(insertUser)
	if err != nil {
		return nil, fmt.Errorf("error when preparing insert statement for user in database, %v", err)
	}
	defer insertStmt.Close()

	if _, err = insertStmt.Exec(user.Username, user.Email, user.Password, user.USD); err != nil {
		if err.(*mysql.MySQLError).Number == 1062 {
			return nil, nil
		}
		return nil, fmt.Errorf("error when inserting user in database, %v", err)
	}
	user.Password = ""
	user.Assets = []model.UserAsset{}
	return &user, nil
}

// Gets all users.
// Returns error on database query error
func (u UsersDBHandler) GetAll() ([]model.User, error) {
	rows, err := u.conn.Query(selectUserAndAssets)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve users from database, %v", err)
	}
	return deserializeUsers(rows)
}

// Gets user with user assets array filled.
// Returns nil user if user does not exist
// Returns error on database query error
func (u UsersDBHandler) GetByUsernameWithAssets(username string) (*model.User, error) {
	rows, err := u.conn.Query(selectUserAndAssetsByUsername, username)
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
		if err := rows.Scan(&asset.user.Username, &asset.user.Email, &asset.user.USD, &asset.assetId, &asset.name, &asset.quantity); err != nil {
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

// Gets user without user assets information.
// Returns nil user if not exists
// Returns error on database query error
func (u UsersDBHandler) GetByUsername(username string) (*model.User, error) {
	row := u.conn.QueryRow(selectUser, username)

	var user model.User
	if err := row.Scan(&user.Username, &user.Email, &user.USD); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not read user row, %v", err)
	}

	return &user, nil
}

// Updates usd value of a user.
// Returns error on database query error
func (u UsersDBHandler) UpdateUSD(username string, money float64) error {
	updateStmt, err := u.conn.Prepare(updateUserUSD)
	if err != nil {
		return fmt.Errorf("error when preparing update statement for user in database, %v", err)
	}
	defer updateStmt.Close()

	res, err := updateStmt.Exec(money, username)
	if err != nil {
		return fmt.Errorf("error when updating user money in database, %v", err)
	}
	cnt, _ := res.RowsAffected()
	if cnt == 0 {
		return fmt.Errorf("could not update the money of user username=%s, usd=%f", username, money)
	}
	return nil
}

// Checks if user with this username and password exists in the database
// Returns error on database query error
func (u UsersDBHandler) Exists(username, password string) (bool, error) {
	row := u.conn.QueryRow(existsUser, username, password)

	var exists int
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("could not check user row existence, %v", err)
	}

	if exists == 1 {
		return true, nil
	}
	return false, nil
}
