package certmanager

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

// Default values for development certificates.
const (
	DefaultOrgName       = "KittyProtocol Dev Environment"
	DefaultCertValidity  = 365 * 24 * time.Hour
	DefaultServerDNSName = "kitty-hub"
)

// SetupTLSConfig loads certificates from disk or generates new self-signed ones.
// This function is intended for development and testing environments.
func SetupTLSConfig(certPath, keyPath string) (*tls.Config, error) {
	// Generate certificates if missing.
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		fmt.Println("[CertManager] No TLS certificates found. Generating new self-signed certificates...")
		if err := GenerateSelfSignedCert(certPath, keyPath, DefaultServerDNSName); err != nil {
			return nil, fmt.Errorf("failed to generate certificates: %w", err)
		}
	}

	// Load certificate pair.
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate files: %w", err)
	}

	// Configure TLS 1.3 with ALPN for QUIC.
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
		NextProtos:   []string{"kitty-quic-v1"},
	}, nil
}

// GenerateSelfSignedCert creates a self-signed ECDSA certificate for development.
func GenerateSelfSignedCert(certPath, keyPath, dnsName string) error {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(DefaultCertValidity)

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{DefaultOrgName},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{dnsName},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	// Ensure directory exists.
	if err := os.MkdirAll(filepath.Dir(certPath), 0755); err != nil {
		return err
	}

	// Write certificate.
	if err := writePEM(certPath, "CERTIFICATE", certDER); err != nil {
		return err
	}

	// Write private key.
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return err
	}
	return writePEM(keyPath, "EC PRIVATE KEY", privBytes)
}

func writePEM(path, pemType string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return pem.Encode(f, &pem.Block{Type: pemType, Bytes: data})
}
