package certmanager

import (
	"crypto/tls"
	"os"
	"testing"
)

// TestSetupTLSConfig verifies that certificates are generated when missing
// and correctly loaded when already present.
func TestSetupTLSConfig(t *testing.T) {
	tempCert := "temp_cert.pem"
	tempKey := "temp_key.pem"

	// Cleanup after test
	defer os.Remove(tempCert)
	defer os.Remove(tempKey)

	// Step 1: Generate new certificates
	cfg, err := SetupTLSConfig(tempCert, tempKey)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if cfg.MinVersion != tls.VersionTLS13 {
		t.Errorf("expected TLS 1.3, got: %x", cfg.MinVersion)
	}

	if _, err := os.Stat(tempCert); os.IsNotExist(err) {
		t.Errorf("certificate file was not generated")
	}

	// Step 2: Load existing certificates
	cfg2, err := SetupTLSConfig(tempCert, tempKey)
	if err != nil {
		t.Fatalf("expected success when loading existing files, got: %v", err)
	}

	if len(cfg2.Certificates) == 0 {
		t.Errorf("expected loaded certificates, got none")
	}
}
