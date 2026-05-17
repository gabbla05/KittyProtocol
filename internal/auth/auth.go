package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var mockUsers = map[string]string{
	"alice": "$2a$10$75AE7Wefqtm/ezWhCP3YR.vooaYVcv6nK/Drn4pK.YH0BbSB.JRPa",
	// Skopiowaliśmy hash od Alice, teraz Bob też ma hasło "secret"
	"bob":     "$2a$10$75AE7Wefqtm/ezWhCP3YR.vooaYVcv6nK/Drn4pK.YH0BbSB.JRPa",
	"charlie": "$2a$10$E1bp02NhmgFw2DR0.l1YruHEbETPvri1nGqBTjY4c9/aLPxrm0Uty",
}

func CheckCredentials(user, pass string) bool {
	fmt.Printf("[AUTH] user=%q pass=%q\n", user, pass)

	hash, exists := mockUsers[user]
	if !exists {
		fmt.Println("[AUTH] user not found")
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		fmt.Println("[AUTH] bcrypt error:", err)
		return false
	}

	fmt.Println("[AUTH] success")
	return true
}

/* Idle Timeout (60s): Mechanizm, który monitoruje każdą sesję. Jeśli przez 60 sekund od klienta nie nadejdzie żadna ramka DATA ani PING, Hub musi jednostronnie zamknąć to połączenie.Rate Limiting (Token Bucket): Musisz zaimplementować faktyczną logikę ograniczania ruchu (10 wiadomości na sekundę na użytkownika / 100 na minutę z IP). Sam komentarz TODO nie obroni serwera przed spamem.Czyszczenie mapy activeSessions: Musisz dopilnować, by w przypadku jakiegokolwiek błędu lub timeoutu, użytkownik był usuwany z mapy w RAM, żeby nie wyciekała pamięć. */
