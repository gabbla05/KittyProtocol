package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

// buildTLSConfig returns a minimal TLS 1.3 configuration.
// Certificate verification is handled manually via TOFU.
func buildTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true, // TOFU: we verify manually
		NextProtos:         []string{"kitty-quic-v1"},
		MinVersion:         tls.VersionTLS13,
	}
}

// verifyOrStoreServerCert performs TOFU:
// - if no pinned cert exists → store serverCert.Raw
// - if pinned cert exists → compare DER bytes
func verifyOrStoreServerCert(cert *x509.Certificate) error {
	if cert == nil {
		return errors.New("no server certificate presented")
	}

	serverDER := cert.Raw

	// First connection → store pinned cert
	if _, err := os.Stat("certs/trusted_cert.pem"); os.IsNotExist(err) {
		pemData := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: serverDER,
		})
		return os.WriteFile("certs/trusted_cert.pem", pemData, 0644)
	}

	// Compare with pinned cert
	trustedPEM, err := os.ReadFile("certs/trusted_cert.pem")
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
