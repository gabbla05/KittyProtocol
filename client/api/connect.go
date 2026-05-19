package api

import (
	"context"
	"errors"

	"github.com/quic-go/quic-go"
)

// Connect establishes a QUIC connection to the Hub and opens a bidirectional stream.
// It performs:
// - TLS 1.3 setup
// - QUIC Dial
// - TOFU certificate verification
// - opening a stream
// - sending HELLO
func (c *KittyClient) Connect(hubAddr string) error {
	if hubAddr == "" {
		return errors.New("hub address is empty")
	}

	tlsConf := buildTLSConfig()

	// QUIC Dial
	conn, err := quic.DialAddr(context.Background(), hubAddr, tlsConf, nil)
	if err != nil {
		return err
	}

	// TOFU certificate verification
	state := conn.ConnectionState()
	if len(state.TLS.PeerCertificates) == 0 {
		conn.CloseWithError(0, "no server certificate")
		return errors.New("no server certificate")
	}

	serverCert := state.TLS.PeerCertificates[0]
	if err := verifyOrStoreServerCert(serverCert); err != nil {
		conn.CloseWithError(0, "certificate verification failed")
		return err
	}

	// Open QUIC stream
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		conn.CloseWithError(0, "stream open failed")
		return err
	}

	// Save connection + stream
	c.mu.Lock()
	c.conn = conn
	c.stream = stream
	c.mu.Unlock()

	// Send HELLO immediately
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

// WaitForHelloOK waits for MEOW_OK after HELLO.
func (c *KittyClient) WaitForHelloOK() error {
	ok, code := c.waitHelloOK()
	if ok {
		c.setState(StateAuthenticating)
		return nil
	}
	return errors.New(code)
}

// WaitForAuthOK waits for MEOW_OK after AUTH.
func (c *KittyClient) WaitForAuthOK() error {
	ok, code := c.waitAuthOK()
	if ok {
		c.setState(StateSelectingTarget)
		return nil
	}
	return errors.New(code)
}
