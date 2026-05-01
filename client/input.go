package main

import (
	"bufio"
	"fmt"
	"strings"
)

// readCredentials reads username and password from stdin.
func readCredentials(r *bufio.Reader) (string, string) {
	fmt.Print("Login: ")
	user, _ := r.ReadString('\n')
	fmt.Print("Hasło: ")
	pass, _ := r.ReadString('\n')
	return strings.TrimSpace(user), strings.TrimSpace(pass)
}

// readTarget reads the target username from stdin.
func readTarget(r *bufio.Reader) string {
	fmt.Print("Do kogo piszesz?: ")
	target, _ := r.ReadString('\n')
	return target[:len(target)-1]
}
