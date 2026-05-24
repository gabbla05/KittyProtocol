package auth

import "testing"

// TestValidateUsername ensures that username validation enforces the expected policy.
func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid_simple", "alice", false},
		{"valid_with_digits", "user123", false},
		{"valid_with_underscore", "user_name", false},
		{"too_short", "ab", true},
		{"too_long", "thisusernameiswaytoolongtobevalid_123", true},
		{"invalid_chars_upper", "Alice", true},
		{"invalid_chars_dash", "user-name", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateUsername(tc.input)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for username %q, got nil", tc.input)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error for username %q: %v", tc.input, err)
			}
		})
	}
}

// TestValidatePassword ensures that the minimum length requirement is enforced.
func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"too_short", "1234567", true},
		{"exact_min", "12345678", false},
		{"longer", "Password123!", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePassword(tc.input)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for password %q, got nil", tc.input)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error for password %q: %v", tc.input, err)
			}
		})
	}
}
