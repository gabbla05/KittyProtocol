package client

import "github.com/joho/godotenv"

// LoadEnv loads .env if present.
// Missing .env is not an error.
func LoadEnv() {
	_ = godotenv.Load()
}
