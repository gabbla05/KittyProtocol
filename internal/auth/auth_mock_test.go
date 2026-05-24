package auth

import "testing"

// TestMockAuthRegisterAndLogin verifies that MockAuth can register and authenticate users.
func TestMockAuthRegisterAndLogin(t *testing.T) {
	m := NewMockAuth()

	user := "testuser"
	pass := "StrongPass123!"

	if err := m.Register(user, pass); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	ok := m.CheckCredentials(user, pass)
	if !ok {
		t.Fatalf("expected credentials to be valid after registration")
	}

	ok = m.CheckCredentials(user, "wrongpass")
	if ok {
		t.Fatalf("expected invalid credentials for wrong password")
	}
}

// TestMockAuthUserExists verifies that UserExists reflects registration state.
func TestMockAuthUserExists(t *testing.T) {
	m := NewMockAuth()

	exists, err := m.UserExists("alice")
	if err != nil {
		t.Fatalf("UserExists returned error: %v", err)
	}
	if !exists {
		t.Fatalf("expected alice to exist in default mock users")
	}

	exists, err = m.UserExists("nonexistent")
	if err != nil {
		t.Fatalf("UserExists returned error: %v", err)
	}
	if exists {
		t.Fatalf("expected nonexistent user to not exist")
	}
}
