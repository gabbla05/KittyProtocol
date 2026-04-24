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


/* Idle Timeout (60s): Mechanizm, który monitoruje każdą sesję. Jeśli przez 60 sekund od klienta nie nadejdzie żadna ramka DATA ani PING, Hub musi jednostronnie zamknąć to połączenie.Rate Limiting (Token Bucket): Musisz zaimplementować faktyczną logikę ograniczania ruchu (10 wiadomości na sekundę na użytkownika / 100 na minutę z IP). Sam komentarz TODO nie obroni serwera przed spamem.Czyszczenie mapy activeSessions: Musisz dopilnować, by w przypadku jakiegokolwiek błędu lub timeoutu, użytkownik był usuwany z mapy w RAM, żeby nie wyciekała pamięć. */