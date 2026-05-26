package auth

// AuthProvider defines the interface for authentication backends.
// Implementations include MockAuth (development) and DBAuth (production).
//
// This interface is ideal as a boundary for higher-level services or
// frontend bindings (e.g. Wails), which can depend on AuthProvider
// without knowing the underlying storage.
type AuthProvider interface {
	// CheckCredentials verifies username and password.
	// Returns true if credentials are valid, false otherwise.
	CheckCredentials(user, pass string) bool

	// Register creates a new user with the given credentials.
	// Implementations MUST:
	//   - validate username and password,
	//   - return an error if the user already exists,
	//   - never log the password.
	Register(user, pass string) error

	// UserExists returns true if the user already exists.
	UserExists(user string) (bool, error)
}
