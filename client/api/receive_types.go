package api

import "encoding/json"

// chatFrameProbe is a lightweight probe structure used to detect
// whether a decrypted DATA payload is a chat control frame.
type chatFrameProbe struct {
	Type    string          `json:"type"`
	From    string          `json:"from"`
	To      string          `json:"to"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type chatRefusePayload struct {
	Reason string `json:"reason,omitempty"`
}

type chatEndPayload struct {
	Reason string `json:"reason,omitempty"`
}

type textMessagePayload struct {
	Text string `json:"text"`
}
