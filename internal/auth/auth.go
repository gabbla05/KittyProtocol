package auth

import "golang.org/x/crypto/bcrypt"

// Mockowa baza danych użytkowników (Login -> Hash BCrypt)
// Hasła: alice -> "secret", bob -> "password"
var mockUsers = map[string]string{
	"alice": "$2a$10$7v.Z6vX7X7X7X7X7X7X7Xu7X7X7X7X7X7X7X7X7X7X7X7X7X7X7X7",
	"bob":   "$2a$10$8v.Z6vX7X7X7X7X7X7X7Xu7X7X7X7X7X7X7X7X7X7X7X7X7X7X7X7",
}

// CheckCredentials weryfikuje czy użytkownik istnieje i czy hasło pasuje.
func CheckCredentials(user, pass string) bool {
	hash, exists := mockUsers[user]
	if !exists {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}
