package main

import (
	"sync"

	"github.com/quic-go/quic-go"
)

// Definicja dostępnych stanów klienta [cite: 4257, 5337]
type ClientState int

const (
	StateDisconnected    ClientState = iota
	StateHandshaking                 // Po wysłaniu HELLO [cite: 4258]
	StateAuthenticating              // Po wysłaniu AUTH [cite: 4263]
	StateSelectingTarget             // Potrzebny do blokowania napływu wiadomości aż do wybory adresata
	StateEstablished                 // Gotowy do wymiany DATA [cite: 4274]
)

type KittyClient struct {
	Conn    quic.Conn
	Stream  quic.Stream
	State   ClientState
	User    string
	Pending map[int64]chan struct{}
	Mu      sync.Mutex
}

// SetState bezpiecznie zmienia stan i loguje przejście
func (c *KittyClient) SetState(newState ClientState) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.State = newState
}
