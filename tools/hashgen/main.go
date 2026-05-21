package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// passes for users:   alice,    bob,        charlie
	passwords := []string{"secret", "password", "private"}

	for _, p := range passwords {
		hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		fmt.Printf("password=%q hash=%q\n", p, string(hash))
	}
}
