package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var mockUsers = map[string]string{
    "alice": "$2a$10$75AE7Wefqtm/ezWhCP3YR.vooaYVcv6nK/Drn4pK.YH0BbSB.JRPa",
    "bob":   "$2a$10$Cxbp6cMDR5S.xNR90lcbSuljSiMEhnCgTF1UWfYGb5VyqSQUVjVri",
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
