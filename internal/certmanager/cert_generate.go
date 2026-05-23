package certmanager

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// GenerateSelfSignedCert creates a fully valid self-signed ECDSA certificate.
// The certificate includes SAN entries for DNS and localhost IPs.
func GenerateSelfSignedCert(certPath, keyPath, dnsName string) error {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate ECDSA key: %w", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(DefaultCertValidity)

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization: []string{DefaultOrgName},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},

		DNSNames: []string{dnsName},
		IPAddresses: []net.IP{
			net.ParseIP("127.0.0.1"),
			net.ParseIP("::1"),
		},

		SignatureAlgorithm: x509.ECDSAWithSHA256,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(certPath), 0755); err != nil {
		return fmt.Errorf("failed to create certificate directory: %w", err)
	}

	if err := writePEM(certPath, "CERTIFICATE", certDER); err != nil {
		return err
	}

	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return fmt.Errorf("failed to marshal EC private key: %w", err)
	}

	return writePEM(keyPath, "EC PRIVATE KEY", privBytes)
}
