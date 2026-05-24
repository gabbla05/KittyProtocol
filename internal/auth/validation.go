package auth

import (
	"fmt"
	"regexp"
)

var usernameRe = regexp.MustCompile(UsernamePattern)

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
	if len(pass) < MinPasswordLength {
		return fmt.Errorf("password too short: minimum %d characters", MinPasswordLength)
	}
	return nil
}
