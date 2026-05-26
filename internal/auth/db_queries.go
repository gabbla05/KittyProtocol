package auth

import "fmt"

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

func (a *DBAuth) insertUser(user, hash string) error {
	_, err := a.db.Exec(
		"INSERT INTO users (username, password_hash) VALUES ($1, $2)",
		user, hash,
	)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}
