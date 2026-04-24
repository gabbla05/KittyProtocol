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

// SetupTLSConfig ładuje certyfikaty z plików lub generuje nowe samopodpisane, jeśli ich brakuje.
func SetupTLSConfig(certPath, keyPath string) (*tls.Config, error) {
	// Sprawdzenie, czy certyfikaty już istnieją na dysku
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		fmt.Println("[CertManager] Brak certyfikatów TLS. Generowanie nowych dla środowiska testowego...")
		if err := generateSelfSignedCert(certPath, keyPath); err != nil {
			return nil, fmt.Errorf("nie udało się wygenerować certyfikatów: %w", err)
		}
	}

	// Wczytanie pary kluczy
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("nie udało się wczytać plików certyfikatu: %w", err)
	}

	// Konfiguracja i wymuszenie TLS 1.3
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
		NextProtos:   []string{"kitty-quic-v1"}, // Wymóg standardu QUIC (ALPN)
	}, nil
}

// generateSelfSignedCert generuje bezpieczny certyfikat ECDSA z krzywą P-256
func generateSelfSignedCert(certPath, keyPath string) error {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // Ważny przez rok

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"KittyProtocol Dev Environment"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	// Upewnienie się, że folder nadrzędny istnieje (np. "certs/")
	os.MkdirAll(filepath.Dir(certPath), 0755)

	// Zapis certyfikatu (.pem)
	certFile, err := os.Create(certPath)
	if err != nil {
		return err
	}
	defer certFile.Close()
	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return err
	}

	// Zapis klucza prywatnego (.pem)
	keyFile, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer keyFile.Close()
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return err
	}
	if err := pem.Encode(keyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}); err != nil {
		return err
	}

	return nil
}
