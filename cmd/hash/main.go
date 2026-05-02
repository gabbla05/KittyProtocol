package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// This utility generates bcrypt hashes for the passwords "secret" and "password". Can be delted after use.
func main() {
	passwords := []string{"secret", "password", "private"}

	for _, p := range passwords {
		hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		fmt.Printf("password=%q hash=%q\n", p, string(hash))
	}
}
