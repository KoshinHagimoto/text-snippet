package dao

import (
	"database/sql"
	"fmt"
	"text-snippet/app/object"
)

type UserDAO struct {
	db *sql.DB
}

func NewUserDAO(db *sql.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (dao *UserDAO) CreateUser(user *object.User) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`
	_, err := dao.db.Exec(query, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	return nil
}

func (dao *UserDAO) GetUserByUsername(username string) (*object.User, error) {
	query := `SELECT id, username, email, password_hash, created_at, email_verified FROM users WHERE username = ?`
	row := dao.db.QueryRow(query, username)

	var user object.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.EmailVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %v", err)
		}
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return &user, nil
}

func (dao *UserDAO) UpdateEmailVerified(userId int, verified bool) error {
	query := `UPDATE users SET email_verified = ? WHERE id = ?`
	_, err := dao.db.Exec(query, verified, userId)
	if err != nil {
		return fmt.Errorf("error updating email verified: %v", err)
	}
	return nil
}
