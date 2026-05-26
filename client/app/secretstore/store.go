package secretstore

import (
	"os"
	"path/filepath"
	"sync"
)

// SecretStore manages per‑peer shared secrets persisted on disk.
// It is safe for concurrent use.
type SecretStore struct {
	mu        sync.Mutex
	path      string
	masterKey []byte
	secrets   map[string][]byte
}

// NewSecretStore creates a SecretStore bound to the given file path.
// If the file exists, it is decrypted; otherwise an empty store is created.
func NewSecretStore(path string, masterKey []byte) *SecretStore {
	s := &SecretStore{
		path:      path,
		masterKey: append([]byte(nil), masterKey...), // defensive copy
		secrets:   make(map[string][]byte),
	}
	_ = s.load()
	return s
}

// Get returns a copy of the secret for a given peer.
func (s *SecretStore) Get(peer string) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	secret, ok := s.secrets[peer]
	if !ok {
		return nil, false
	}
	out := make([]byte, len(secret))
	copy(out, secret)
	return out, true
}

// Set stores/updates the secret for a given peer and persists the file.
func (s *SecretStore) Set(peer string, secret []byte) error {
	if peer == "" {
		return ErrEmptyPeer
	}
	if len(secret) == 0 {
		return ErrEmptySecret
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	buf := make([]byte, len(secret))
	copy(buf, secret)
	s.secrets[peer] = buf

	return s.saveLocked()
}

// All returns a deep copy of all secrets.
func (s *SecretStore) All() map[string][]byte {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make(map[string][]byte, len(s.secrets))
	for k, v := range s.secrets {
		buf := make([]byte, len(v))
		copy(buf, v)
		out[k] = buf
	}
	return out
}

// load decrypts and loads secrets from disk.
func (s *SecretStore) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	plaintext, err := decrypt(s.masterKey, string(data))
	if err != nil {
		return err
	}

	m, err := unmarshalDisk(plaintext)
	if err != nil {
		return err
	}

	s.secrets = m
	return nil
}

// saveLocked serializes and encrypts the current secrets map to disk.
// Caller must hold s.mu.
func (s *SecretStore) saveLocked() error {
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	plaintext, err := marshalDisk(s.secrets)
	if err != nil {
		return err
	}

	enc, err := encrypt(s.masterKey, plaintext)
	if err != nil {
		return err
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, []byte(enc), 0o600); err != nil {
		return err
	}

	return os.Rename(tmp, s.path)
}
