package api

import (
	"context"
	"errors"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/quic-go/quic-go"
)

// Connect establishes a fresh QUIC connection to the Hub and opens a bidirectional stream.
//
// BEHAVIOR:
//   - If the client was previously connected, Connect() implicitly closes the old connection.
//   - Performs TLS 1.3 setup, QUIC dial, TOFU certificate verification, and stream opening.
//   - Sends the HELLO frame immediately after the stream is opened (non‑blocking).
//   - Does NOT start receiver or ping loops — the caller must start them manually.
//
// STATE TRANSITIONS:
//
//	StateDisconnected → StateHandshaking
func (c *KittyClient) Connect(hubAddr string) error {
	if hubAddr == "" {
		return errors.New("hub address is empty")
	}

	// Ensure clean state if reconnecting
	c.Close()

	// Recreate control channels for background loops.
	c.mu.Lock()
	c.stopPing = make(chan struct{})
	c.stopRecv = make(chan struct{})
	c.mu.Unlock()

	tlsConf := buildTLSConfig()

	// QUIC Dial (real transport)
	rawConn, err := quic.DialAddr(context.Background(), hubAddr, tlsConf, nil)
	if err != nil {
		return err
	}
	conn := newQuicConnAdapter(rawConn)

	// TOFU certificate verification
	state := conn.ConnectionState()
	if len(state.TLS.PeerCertificates) == 0 {
		_ = conn.CloseWithError(0, "no server certificate")
		return errors.New("no server certificate")
	}

	serverCert := state.TLS.PeerCertificates[0]
	if err := verifyOrStoreServerCert(serverCert); err != nil {
		_ = conn.CloseWithError(0, "certificate verification failed")
		return err
	}

	// Open QUIC stream via adapter
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		_ = conn.CloseWithError(0, "stream open failed")
		return err
	}

	// Save connection + stream
	c.mu.Lock()
	c.conn = conn
	c.stream = stream

	// Reset session‑level state
	c.target = ""
	c.lastFrame = nil
	c.replay = protection.NewReplayDetector()
	c.ackMgr = NewAckManager()
	c.mu.Unlock()

	// Send HELLO immediately (async handshake)
	if err := c.SendHello(); err != nil {
		return err
	}

	c.setState(StateHandshaking)
	return nil
}

// Disconnect closes the QUIC connection and stream.
func (c *KittyClient) Disconnect() {
	c.Close()
}
