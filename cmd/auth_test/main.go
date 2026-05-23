package main

import (
	"database/sql"
	"fmt"

	"github.com/gabbla05/KittyProtocol/internal/auth"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://kitty:kittypass@localhost:5432/kittyhub?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	a := auth.NewDBAuth(db)

	// Register
	if err := a.Register("testuser", "supersecret"); err != nil {
		fmt.Println("Register error:", err)
	} else {
		fmt.Println("Register OK")
	}

	// CheckCredentials (OK)
	if a.CheckCredentials("testuser", "supersecret") {
		fmt.Println("Login OK")
	} else {
		fmt.Println("Login FAILED")
	}

	// CheckCredentials (bad password)
	if a.CheckCredentials("testuser", "wrongpass") {
		fmt.Println("Login should NOT succeed with wrong password")
	} else {
		fmt.Println("Login failed as expected (wrong password)")
	}
}
