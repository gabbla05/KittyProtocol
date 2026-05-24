package auth

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
)

// TestDBAuthUserExists verifies that UserExists queries the database correctly.
func TestDBAuthUserExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	a := NewDBAuth(db)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)",
	)).
		WithArgs("alice").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := a.UserExists("alice")
	if err != nil {
		t.Fatalf("UserExists returned error: %v", err)
	}
	if !exists {
		t.Fatalf("expected alice to exist")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// TestDBAuthCheckCredentials verifies that CheckCredentials compares bcrypt hashes.
func TestDBAuthCheckCredentials(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	a := NewDBAuth(db)

	pass := "Secret123!"
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT password_hash FROM users WHERE username=$1",
	)).
		WithArgs("alice").
		WillReturnRows(sqlmock.NewRows([]string{"password_hash"}).AddRow(string(hash)))

	ok := a.CheckCredentials("alice", pass)
	if !ok {
		t.Fatalf("expected valid credentials to pass")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
