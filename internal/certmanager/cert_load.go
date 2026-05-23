package certmanager

import (
	"crypto/tls"
	"fmt"
	"os"
)

// SetupTLSConfig loads existing certificates or generates new self-signed ones.
// This function is used by the Hub during startup.
func SetupTLSConfig(certPath, keyPath string) (*tls.Config, error) {
	// Generate certificates if missing
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		fmt.Println("[CertManager] No TLS certificates found. Generating new self-signed certificates...")
		if err := GenerateSelfSignedCert(certPath, keyPath, DefaultServerDNSName); err != nil {
			return nil, fmt.Errorf("failed to generate certificates: %w", err)
		}
	}

	// Load certificate pair
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate files: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
		NextProtos:   []string{DefaultALPNProtocol},
	}, nil
}
