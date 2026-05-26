package api

// OpResult represents the outcome of an asynchronous operation such as
// HELLO, AUTH or REGISTER. It implements the error interface for convenient
// propagation through Go APIs.
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

// ChatRequestEvent is emitted when a peer requests to start a chat session.
type ChatRequestEvent struct {
	From string
}

// ChatAcceptEvent is emitted when a peer accepts a previously sent chat request.
type ChatAcceptEvent struct {
	From string
}

// ChatRefuseEvent is emitted when a peer refuses a chat request.
type ChatRefuseEvent struct {
	From   string
	Reason string
}

// ChatEndEvent is emitted when a peer terminates an active chat session.
type ChatEndEvent struct {
	From   string
	Reason string
}

// ChatMessageEvent is emitted when a peer sends a text message within
// an established chat session.
type ChatMessageEvent struct {
	From string
	Text string
}

// ChatRequestEvents returns a read-only channel of incoming chat requests.
func (c *KittyClient) ChatRequestEvents() <-chan ChatRequestEvent {
	return c.chatReqCh
}

// ChatAcceptEvents returns a read-only channel of chat accept events.
func (c *KittyClient) ChatAcceptEvents() <-chan ChatAcceptEvent {
	return c.chatAcceptCh
}

// ChatRefuseEvents returns a read-only channel of chat refusal events.
func (c *KittyClient) ChatRefuseEvents() <-chan ChatRefuseEvent {
	return c.chatRefuseCh
}

// ChatEndEvents returns a read-only channel of chat termination events.
func (c *KittyClient) ChatEndEvents() <-chan ChatEndEvent {
	return c.chatEndCh
}

// ChatMessageEvents returns a read-only channel of incoming chat messages.
func (c *KittyClient) ChatMessageEvents() <-chan ChatMessageEvent {
	return c.chatMsgCh
}
