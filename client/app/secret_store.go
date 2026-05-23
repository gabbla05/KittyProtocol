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
// Each Kitty user has its own directory: ~/.kitty/<kittyUser>/secrets.json
type SecretStore struct {
	mu      sync.Mutex
	path    string
	secrets map[string][]byte
}

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

// PathForUser returns ~/.kitty/<kittyUser>/secrets.json
func PathForUser(kittyUser string) string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(".", "kitty", kittyUser, "secrets.json")
	}
	return filepath.Join(home, ".kitty", kittyUser, "secrets.json")
}

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

func (s *SecretStore) Set(peer string, secret []byte) error {
	if peer == "" {
		return errors.New("peer cannot be empty")
	}
	if len(secret) == 0 {
		return errors.New("secret cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	buf := make([]byte, len(secret))
	copy(buf, secret)
	s.secrets[peer] = buf

	return s.saveLocked()
}

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

	var ds diskSecrets
	if err := json.Unmarshal(data, &ds); err != nil {
		return err
	}

	s.secrets = make(map[string][]byte, len(ds.Peers))
	for peer, b64 := range ds.Peers {
		raw, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			continue
		}
		s.secrets[peer] = raw
	}

	return nil
}

func (s *SecretStore) saveLocked() error {
	dir := filepath.Dir(s.path)

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

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}

	return os.Rename(tmp, s.path)
}

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
