package auth

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// DBAuth is a PostgreSQL-backed authentication provider.
type DBAuth struct {
	db *sql.DB
}

func NewDBAuth(db *sql.DB) *DBAuth {
	return &DBAuth{db: db}
}

func (a *DBAuth) CheckCredentials(user, pass string) bool {
	hash, err := a.lookupPasswordHash(user)
	if err != nil {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
}

func (a *DBAuth) Register(user, pass string) error {
	if err := validateUsername(user); err != nil {
		return err
	}
	if err := validatePassword(pass); err != nil {
		return err
	}

	exists, err := a.UserExists(user)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return fmt.Errorf("username already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return a.insertUser(user, string(hash))
}
