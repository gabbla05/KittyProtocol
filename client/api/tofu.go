// tofu.go
// Implements Trust-On-First-Use (TOFU) certificate pinning for KittyClient.
// This is an application-level trust mechanism independent of TLS transport.

package api

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
)

const (
	// pinnedCertDir is the directory where the pinned certificate is stored.
	pinnedCertDir  = "certs"
	pinnedCertFile = "trusted_cert.pem"
)

// verifyOrStoreServerCert performs TOFU certificate pinning.
//
// STORAGE:
//   - The pinned certificate is stored as PEM at "certs/trusted_cert.pem"
//     relative to the current working directory.
//
// BEHAVIOR:
//   - If no pinned certificate exists → store the presented certificate.
//   - If a pinned certificate exists → compare DER bytes.
//   - Any mismatch results in an error (possible MITM).
func verifyOrStoreServerCert(cert *x509.Certificate) error {
	if cert == nil {
		return errors.New("no server certificate presented")
	}

	serverDER := cert.Raw
	pinnedPath := filepath.Join(pinnedCertDir, pinnedCertFile)

	// First run: no pinned certificate → store current one.
	if _, err := os.Stat(pinnedPath); os.IsNotExist(err) {
		if err := os.MkdirAll(pinnedCertDir, 0o755); err != nil {
			return err
		}

		pemData := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: serverDER,
		})

		// 0600: only current user can read/write pinned cert.
		return os.WriteFile(pinnedPath, pemData, 0o600)
	}

	// Subsequent runs: compare with pinned certificate.
	trustedPEM, err := os.ReadFile(pinnedPath)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(trustedPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return errors.New("trusted_cert.pem is invalid")
	}

	if !bytes.Equal(block.Bytes, serverDER) {
		return errors.New("server certificate mismatch (possible MITM)")
	}

	return nil
}
