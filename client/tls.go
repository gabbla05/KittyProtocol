// client/tls.go
package main

import "crypto/tls"

func buildTLSConfig() (*tls.Config, error) {
	return &tls.Config{
		InsecureSkipVerify: true, // pozwalamy na dowolny cert, ale sami go sprawdzimy (TOFU)
		NextProtos:         []string{"kitty-quic-v1"},
		MinVersion:         tls.VersionTLS13,
	}, nil
}
