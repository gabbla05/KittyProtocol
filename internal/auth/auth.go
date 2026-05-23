package auth

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// AuthProvider defines the interface for authentication backends.
// This allows swapping mock auth for a real database implementation.
type AuthProvider interface {
	// CheckCredentials verifies username and password.
	// Returns true if credentials are valid, false otherwise.
	CheckCredentials(user, pass string) bool

	// Register creates a new user with the given credentials.
	// Implementations MUST:
	//   - validate username and password,
	//   - return an error if the user already exists,
	//   - never log the password.
	Register(user, pass string) error

	// UserExists returns true if the user already exists.
	UserExists(user string) (bool, error)
}

// -----------------------------------------------------------------------------
// MockAuth – in‑memory implementation for development/testing
// -----------------------------------------------------------------------------

// MockAuth is a simple in-memory authentication provider.
// Intended ONLY for development and testing.
type MockAuth struct {
	users map[string]string // username -> bcrypt hash
}

// NewMockAuth creates a mock authentication provider with predefined users.
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

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
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

// -----------------------------------------------------------------------------
// Shared validation helpers
// -----------------------------------------------------------------------------

var usernameRe = regexp.MustCompile(`^[a-z0-9_]{3,32}$`)

// validateUsername enforces a simple, predictable username policy.
func validateUsername(user string) error {
	if !usernameRe.MatchString(user) {
		return fmt.Errorf("invalid username: must be 3–32 chars, [a-z0-9_]")
	}
	return nil
}

// validatePassword enforces a minimal password policy.
// You can tighten this later (e.g. require digits/symbols).
func validatePassword(pass string) error {
	if len(pass) < 8 {
		return fmt.Errorf("password too short: minimum 8 characters")
	}
	return nil
}
