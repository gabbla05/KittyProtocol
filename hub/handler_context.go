package main

import (
	"fmt"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/quic-go/quic-go"
)

type connectionState int

const (
	stateInit connectionState = iota
	stateHelloReceived
	stateAuthenticated
)

type clientContext struct {
	conn      *quic.Conn
	stream    *quic.Stream
	session   *protection.Session
	username  string
	authTimer *protection.AuthTimer
	state     connectionState
}

func (c *clientContext) cleanup() {
	if c.session != nil {
		fmt.Println("[Handler: Context] Cleaning up session for:", c.username)
		globalSessions.Remove(c.username)
		if c.session.CloseFunc != nil {
			c.session.CloseFunc()
		}
		c.session = nil
	}
	if c.authTimer != nil {
		c.authTimer.Stop()
		c.authTimer = nil
	}
}

func (c *clientContext) touch() {
	if c.session != nil {
		c.session.LastActive = time.Now()
	}
}
