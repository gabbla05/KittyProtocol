package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// MockAuth is a simple in-memory authentication provider.
// Intended ONLY for development and testing. It is not persistent.
type MockAuth struct {
	users map[string]string // username -> bcrypt hash
}

// NewMockAuth creates a mock authentication provider with predefined users.
// The initial users are intended for local development and manual testing.
func NewMockAuth() *MockAuth {
	return &MockAuth{
		users: map[string]string{
			"alice":   "$2a$10$75AE7Wefqtm/ezWhCP3YR.vooaYVcv6nK/Drn4pK.YH0BbSB.JRPa",
			"bob":     "$2a$10$Cxbp6cMDR5S.xNR90lcbSuljSiMEhnCgTF1UWfYGb5VyqSQUVjVri",
			"charlie": "$2a$10$E1bp02NhmgFw2DR0.l1YruHEbETPvri1nGqBTjY4c9/aLPxrm0Uty",
		},
	}
}

// CheckCredentials verifies username and password using bcrypt.
// NOTE: Passwords are NEVER logged for security reasons.
func (m *MockAuth) CheckCredentials(user, pass string) bool {
	fmt.Printf("[AUTH] user=%q\n", user)

	hash, exists := m.users[user]
	if !exists {
		fmt.Println("[AUTH] user not found")
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)); err != nil {
		fmt.Println("[AUTH] invalid password")
		return false
	}

	fmt.Println("[AUTH] success")
	return true
}

// Register adds a new user to the in-memory store.
// This is only for development/testing; no persistence is performed.
func (m *MockAuth) Register(user, pass string) error {
	if err := validateUsername(user); err != nil {
		return err
	}
	if err := validatePassword(pass); err != nil {
		return err
	}

	if _, exists := m.users[user]; exists {
		return fmt.Errorf("username already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	m.users[user] = string(hash)
	return nil
}

// UserExists checks if the user is present in the in-memory store.
func (m *MockAuth) UserExists(user string) (bool, error) {
	_, exists := m.users[user]
	return exists, nil
}
