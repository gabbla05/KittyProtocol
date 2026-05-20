package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

// buildTLSConfig returns a minimal TLS 1.3 configuration for QUIC.
// Certificate verification is intentionally disabled (InsecureSkipVerify)
// because KittyClient performs TOFU (Trust On First Use) manually.
//
// SECURITY MODEL:
//   - First connection: the server certificate is stored locally.
//   - Subsequent connections: the certificate must match the stored one.
//   - Any mismatch is treated as a potential MITM attack.
func buildTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true, // TOFU: manual verification
		NextProtos:         []string{"kitty-quic-v1"},
		MinVersion:         tls.VersionTLS13,
	}
}

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
	pinnedPath := "certs/trusted_cert.pem"

	if _, err := os.Stat(pinnedPath); os.IsNotExist(err) {
		pemData := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: serverDER,
		})
		return os.WriteFile(pinnedPath, pemData, 0644)
	}

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
