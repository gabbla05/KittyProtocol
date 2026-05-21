package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// AuthProvider defines the interface for authentication backends.
// This allows swapping mock auth for a real database implementation.
type AuthProvider interface {
	CheckCredentials(user, pass string) bool
}

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
