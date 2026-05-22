package app

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// SecretStore manages per-peer shared secrets persisted on disk.
//
// SECURITY:
//   - Secrets are stored in a file with 0600 permissions inside a directory
//     with 0700 permissions.
//   - Secrets are stored as base64-encoded bytes (no hashing, no key wrapping).
//   - For a production-grade system, secrets should be encrypted at rest
//     (e.g. OS keyring, hardware-backed keystore, or master password).
type SecretStore struct {
	mu      sync.Mutex
	path    string
	secrets map[string][]byte
}

// diskSecrets is the JSON representation persisted on disk.
type diskSecrets struct {
	Peers map[string]string `json:"peers"` // peer -> base64(secret)
}

// NewSecretStore creates a SecretStore bound to the given file path.
// If the file exists, it is loaded; otherwise an empty store is created.
func NewSecretStore(path string) *SecretStore {
	s := &SecretStore{
		path:    path,
		secrets: make(map[string][]byte),
	}
	_ = s.load()
	return s
}

// Get returns the shared secret for a given peer, if present.
func (s *SecretStore) Get(peer string) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	secret, ok := s.secrets[peer]
	if !ok {
		return nil, false
	}
	// Return a copy to avoid accidental mutation.
	out := make([]byte, len(secret))
	copy(out, secret)
	return out, true
}

// Set stores/updates the shared secret for a given peer and persists it to disk.
func (s *SecretStore) Set(peer string, secret []byte) error {
	if peer == "" {
		return errors.New("peer cannot be empty")
	}
	if len(secret) == 0 {
		return errors.New("secret cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Store a copy in memory.
	buf := make([]byte, len(secret))
	copy(buf, secret)
	s.secrets[peer] = buf

	return s.saveLocked()
}

// load reads the secrets file from disk, if it exists.
func (s *SecretStore) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			// No file yet — start with empty store.
			return nil
		}
		return err
	}

	var ds diskSecrets
	if err := json.Unmarshal(data, &ds); err != nil {
		return err
	}

	s.secrets = make(map[string][]byte, len(ds.Peers))
	for peer, b64 := range ds.Peers {
		raw, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			continue // skip malformed entries
		}
		s.secrets[peer] = raw
	}

	return nil
}

// saveLocked writes the current secrets map to disk.
// Caller must hold s.mu.
func (s *SecretStore) saveLocked() error {
	dir := filepath.Dir(s.path)

	// Ensure directory exists with restrictive permissions.
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	ds := diskSecrets{
		Peers: make(map[string]string, len(s.secrets)),
	}
	for peer, secret := range s.secrets {
		ds.Peers[peer] = base64.StdEncoding.EncodeToString(secret)
	}

	data, err := json.MarshalIndent(ds, "", "  ")
	if err != nil {
		return err
	}

	// Write to a temp file and then atomically rename.
	tmpPath := s.path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o600); err != nil {
		return err
	}

	return os.Rename(tmpPath, s.path)
}

// defaultSecretStorePath returns the default path for the secret store file.
func defaultSecretStorePath() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		// Fallback: current working directory.
		return "kitty_secrets.json"
	}
	return filepath.Join(home, ".kitty", "secrets.json")
}
