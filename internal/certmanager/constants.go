package certmanager

import "time"

// DefaultOrgName is used as the Organization field in self-signed certificates.
const DefaultOrgName = "KittyProtocol Dev Environment"

// DefaultCertValidity defines how long generated certificates remain valid.
const DefaultCertValidity = 365 * 24 * time.Hour

// DefaultServerDNSName is the SAN DNS entry used for the Hub.
const DefaultServerDNSName = "kitty-hub"

// DefaultALPNProtocol is the ALPN identifier used for QUIC connections.
const DefaultALPNProtocol = "kitty-quic-v1"
