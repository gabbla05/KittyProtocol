package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// SecretStore manages per-peer shared secrets persisted on disk.
// Each Kitty user has its own directory: ~/.kitty/<kittyUser>/secrets.json.enc
//
// Plik na dysku jest CAŁY zaszyfrowany AES-GCM kluczem wyprowadzonym
// z masterKey (np. hasła użytkownika).
type SecretStore struct {
	mu        sync.Mutex
	path      string
	masterKey []byte
	secrets   map[string][]byte
}

type diskSecrets struct {
	Peers map[string]string `json:"peers"` // peer -> base64(secret)
}

// deriveKey normalizuje masterKey do 32 bajtów (AES-256) przez SHA-256.
// Jeśli masterKey jest hasłem, to jest to prosty KDF.
// W przyszłości można to podmienić na PBKDF2/Argon2.
func deriveKey(masterKey []byte) []byte {
	sum := sha256.Sum256(masterKey)
	return sum[:]
}

// encrypt encryptuje plaintext przy użyciu AES-GCM(masterKey).
// Zwraca: base64( nonce || ciphertext ).
func encrypt(masterKey, plaintext []byte) (string, error) {
	if len(masterKey) == 0 {
		return "", errors.New("master key is empty")
	}

	key := deriveKey(masterKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	out := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(out), nil
}

// decrypt odszyfrowuje base64( nonce || ciphertext ) przy użyciu AES-GCM(masterKey).
func decrypt(masterKey []byte, enc string) ([]byte, error) {
	if len(masterKey) == 0 {
		return nil, errors.New("master key is empty")
	}

	raw, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return nil, err
	}

	key := deriveKey(masterKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(raw) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce := raw[:gcm.NonceSize()]
	ciphertext := raw[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// NewSecretStore creates a SecretStore bound to the given file path.
// masterKey musi być stały dla danego użytkownika (np. hasło logowania).
// Jeśli plik istnieje, jest odszyfrowywany; w przeciwnym razie tworzony jest pusty store.
func NewSecretStore(path string, masterKey []byte) *SecretStore {
	s := &SecretStore{
		path:      path,
		masterKey: append([]byte(nil), masterKey...),
		secrets:   make(map[string][]byte),
	}
	_ = s.load()
	return s
}

// PathForUser returns ~/.kitty/<kittyUser>/secrets.json.enc
func PathForUser(kittyUser string) string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(".", "kitty", kittyUser, "secrets.json.enc")
	}
	return filepath.Join(home, ".kitty", kittyUser, "secrets.json.enc")
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

	// odszyfruj cały plik
	plaintext, err := decrypt(s.masterKey, string(data))
	if err != nil {
		return err
	}

	var ds diskSecrets
	if err := json.Unmarshal(plaintext, &ds); err != nil {
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

	plaintext, err := json.MarshalIndent(ds, "", "  ")
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
