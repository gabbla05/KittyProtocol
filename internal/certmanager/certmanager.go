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
	"net"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultOrgName       = "KittyProtocol Dev Environment"
	DefaultCertValidity  = 365 * 24 * time.Hour
	DefaultServerDNSName = "kitty-hub"
)

// SetupTLSConfig loads certificates from disk or generates new self-signed ones.
func SetupTLSConfig(certPath, keyPath string) (*tls.Config, error) {
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		fmt.Println("[CertManager] No TLS certificates found. Generating new self-signed certificates...")
		if err := GenerateSelfSignedCert(certPath, keyPath, DefaultServerDNSName); err != nil {
			return nil, fmt.Errorf("failed to generate certificates: %w", err)
		}
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate files: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
		NextProtos:   []string{"kitty-quic-v1"},
	}, nil
}

// GenerateSelfSignedCert creates a fully valid self-signed ECDSA certificate.
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
		NotBefore: notBefore,
		NotAfter:  notAfter,

		// Required for self-signed certs
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage: x509.KeyUsageDigitalSignature |
			x509.KeyUsageCertSign,

		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},

		// SAN
		DNSNames: []string{dnsName},
		IPAddresses: []net.IP{
			net.ParseIP("127.0.0.1"),
			net.ParseIP("::1"),
		},

		// Recommended for compatibility
		SignatureAlgorithm: x509.ECDSAWithSHA256,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(certPath), 0755); err != nil {
		return err
	}

	if err := writePEM(certPath, "CERTIFICATE", certDER); err != nil {
		return err
	}

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
