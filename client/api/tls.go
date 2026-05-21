// tls.go
// Hardened TLS 1.3 configuration for QUIC + TOFU verification.

package api

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"time"
)

// buildTLSConfig returns a hardened TLS 1.3 configuration.
//
// SECURITY MODEL
//   - TLS 1.3 only
//   - QUIC-specific ALPN enforced
//   - TOFU certificate validation
//   - Session resumption disabled
//   - Strong curves only
//   - Strong TLS 1.3 cipher suites only
//
// CERTIFICATE VALIDATION
// Default PKI verification is intentionally disabled because
// the application uses TOFU (Trust On First Use).
//
// The first successfully seen certificate is stored locally.
// Future connections must present the same certificate.
func buildTLSConfig() *tls.Config {
	return &tls.Config{
		// TLS 1.3 only
		MinVersion: tls.VersionTLS13,
		MaxVersion: tls.VersionTLS13,

		// QUIC ALPN
		NextProtos: []string{
			"kitty-quic-v1",
		},

		// Disable TLS session resumption.
		// Prevents silently persisting compromised sessions.
		SessionTicketsDisabled: true,

		// Explicitly disable renegotiation.
		// (TLS 1.3 already removed it.)
		Renegotiation: tls.RenegotiateNever,

		// Strong ECDHE curves only
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
		},

		// Explicit TLS 1.3 cipher selection
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},

		// Disable standard PKI verification.
		// We validate manually via TOFU.
		InsecureSkipVerify: true,

		// Modern verification hook.
		VerifyConnection: func(cs tls.ConnectionState) error {
			if len(cs.PeerCertificates) == 0 {
				return errors.New("server did not provide certificate")
			}

			cert := cs.PeerCertificates[0]

			// Basic certificate sanity checks.
			if err := validateServerCertificate(cert); err != nil {
				return err
			}

			// TOFU validation / pinning.
			return verifyOrStoreServerCert(cert)
		},
	}
}

// validateServerCertificate performs minimal cryptographic sanity checks.
//
// NOTE:
// This does NOT perform CA/hostname validation.
// TOFU replaces the traditional PKI trust model.
func validateServerCertificate(cert *x509.Certificate) error {
	if err := cert.CheckSignatureFrom(cert); err != nil {
		return errors.New("invalid self-signed certificate")
	}

	now := time.Now()

	if now.Before(cert.NotBefore) {
		return errors.New("certificate not yet valid")
	}

	if now.After(cert.NotAfter) {
		return errors.New("certificate expired")
	}

	return nil
}
