package api

import (
	"context"
	"sync"

	"github.com/gabbla05/KittyProtocol/internal/protection"
)

type ClientState int

type AppPayloadHandler func(sender string, payload []byte)
type ErrorHandler func(code, desc string)
type StatusHandler func(target, status string)
type DisconnectHandler func(err error)

const (
	StateDisconnected ClientState = iota
	StateHandshaking
	StateAuthenticating
	StateRegistering
	StateSelectingTarget
	StateEstablished
)

type peerKeys struct {
	kEnc []byte
	kMac []byte
}

type KittyClient struct {
	mu    sync.Mutex
	state ClientState

	conn   ConnAdapter
	stream StreamAdapter

	user   string
	target string

	ackMgr   *AckManager
	replay   *protection.ReplayDetector
	stopPing chan struct{}
	stopRecv chan struct{}
	ctx      context.Context
	cancel   context.CancelFunc

	lastFrame []byte

	peerKeys map[string]peerKeys

	appHandler        AppPayloadHandler
	errHandler        ErrorHandler
	statusHandler     StatusHandler
	disconnectHandler DisconnectHandler

	helloCh    chan OpResult
	authCh     chan OpResult
	registerCh chan OpResult

	chatReqCh    chan ChatRequestEvent
	chatAcceptCh chan ChatAcceptEvent
	chatRefuseCh chan ChatRefuseEvent
	chatEndCh    chan ChatEndEvent
	chatMsgCh    chan ChatMessageEvent
}

func NewKittyClient() *KittyClient {
	ctx, cancel := context.WithCancel(context.Background())

	return &KittyClient{
		state:      StateDisconnected,
		ackMgr:     NewAckManager(),
		replay:     protection.NewReplayDetector(),
		stopPing:   make(chan struct{}),
		stopRecv:   make(chan struct{}),
		ctx:        ctx,
		cancel:     cancel,
		peerKeys:   make(map[string]peerKeys),
		helloCh:    make(chan OpResult, 1),
		authCh:     make(chan OpResult, 1),
		registerCh: make(chan OpResult, 1),

		chatReqCh:    make(chan ChatRequestEvent, 4),
		chatAcceptCh: make(chan ChatAcceptEvent, 4),
		chatRefuseCh: make(chan ChatRefuseEvent, 4),
		chatEndCh:    make(chan ChatEndEvent, 4),
		chatMsgCh:    make(chan ChatMessageEvent, 16),
	}
}

func (c *KittyClient) RegisterAppPayloadHandler(h AppPayloadHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.appHandler = h
}

func (c *KittyClient) OnError(h ErrorHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errHandler = h
}

func (c *KittyClient) OnStatus(h StatusHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.statusHandler = h
}

func (c *KittyClient) OnDisconnected(h DisconnectHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.disconnectHandler = h
}

func (c *KittyClient) User() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.user
}

func (c *KittyClient) getKeysForPeer(peer string) (kEnc, kMac []byte, ok bool) {
	pk, exists := c.peerKeys[peer]
	if !exists {
		return nil, nil, false
	}
	return pk.kEnc, pk.kMac, true
}

type OpResult struct {
	OK   bool
	Code string
	Desc string
}

func (r OpResult) Error() string {
	if r.OK {
		return ""
	}
	if r.Desc != "" {
		return r.Code + ": " + r.Desc
	}
	return r.Code
}

type ChatRequestEvent struct {
	From string
}

type ChatAcceptEvent struct {
	From string
}

type ChatRefuseEvent struct {
	From   string
	Reason string
}

type ChatEndEvent struct {
	From   string
	Reason string
}

type ChatMessageEvent struct {
	From string
	Text string
}

// Chat event getters
func (c *KittyClient) ChatRequestEvents() <-chan ChatRequestEvent {
	return c.chatReqCh
}

func (c *KittyClient) ChatAcceptEvents() <-chan ChatAcceptEvent {
	return c.chatAcceptCh
}

func (c *KittyClient) ChatRefuseEvents() <-chan ChatRefuseEvent {
	return c.chatRefuseCh
}

func (c *KittyClient) ChatEndEvents() <-chan ChatEndEvent {
	return c.chatEndCh
}

func (c *KittyClient) ChatMessageEvents() <-chan ChatMessageEvent {
	return c.chatMsgCh
}
