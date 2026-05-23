package protocol

// Frame type identifiers used across the KittyProtocol transport layer.
// These values appear in the "type" field of every JSON frame.
const (
	FrameTypeHello     = "HELLO"
	FrameTypeAuth      = "AUTH"
	FrameTypeRegister  = "REGISTER"
	FrameTypeData      = "DATA"
	FrameTypeMeowOK    = "MEOW_OK"
	FrameTypeError     = "ERROR"
	FrameTypeGetStatus = "GET_STATUS"
	FrameTypeStatusRes = "STATUS_RES"
	FrameTypePing      = "PING"
	FrameTypeBye       = "BYE"
)

// CurrentProtocolVersion defines the version of the KittyProtocol.
// Clients must send this value in the HELLO frame.
const CurrentProtocolVersion = "1.0"
