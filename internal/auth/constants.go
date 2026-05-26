package auth

// MinPasswordLength defines the minimum allowed password length.
// This is enforced by validatePassword.
const MinPasswordLength = 8

// UsernamePattern defines the allowed username format:
//   - 3–32 characters
//   - lowercase letters, digits and underscore only.
const UsernamePattern = `^[a-z0-9_]{3,32}$`
