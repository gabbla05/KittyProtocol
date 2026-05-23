package protocol

// MeowOkFrame is sent by the Hub to acknowledge successful operations
// such as HELLO, AUTH, REGISTER, or other protocol-level actions.
type MeowOkFrame struct {
	BaseFrame
	Status string `json:"status,omitempty"`
}
