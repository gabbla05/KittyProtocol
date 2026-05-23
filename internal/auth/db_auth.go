package auth

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// DBAuth is a PostgreSQL-backed authentication provider.
// It assumes a table:
//
//	CREATE TABLE users (
//	    id            SERIAL PRIMARY KEY,
//	    username      TEXT UNIQUE NOT NULL,
//	    password_hash TEXT NOT NULL,
//	    created_at    TIMESTAMP NOT NULL DEFAULT NOW()
//	);
type DBAuth struct {
	db *sql.DB
}

func NewDBAuth(db *sql.DB) *DBAuth {
	return &DBAuth{db: db}
}

// CheckCredentials verifies username and password against the database.
func (a *DBAuth) CheckCredentials(user, pass string) bool {
	hash, err := a.lookupPasswordHash(user)
	if err != nil {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
}

// Register creates a new user in the database with a bcrypt-hashed password.
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

	_, err = a.db.Exec(
		"INSERT INTO users (username, password_hash) VALUES ($1, $2)",
		user, string(hash),
	)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

// UserExists checks if a user with the given username exists.
func (a *DBAuth) UserExists(user string) (bool, error) {
	var exists bool
	err := a.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)",
		user,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// lookupPasswordHash returns the stored bcrypt hash for a user.
func (a *DBAuth) lookupPasswordHash(user string) (string, error) {
	var hash string
	err := a.db.QueryRow(
		"SELECT password_hash FROM users WHERE username=$1",
		user,
	).Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}
