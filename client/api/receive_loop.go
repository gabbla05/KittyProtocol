package api

import (
	"encoding/json"

	"github.com/gabbla05/KittyProtocol/protocol"
)

// StartReceiverLoop starts a background goroutine that continuously reads
// frames from the QUIC stream, parses them and dispatches to type-specific
// handlers. It closes the provided 'disconnected' channel exactly once when
// the underlying stream is broken or the loop terminates.
func (c *KittyClient) StartReceiverLoop(disconnected chan struct{}) {
	c.mu.Lock()
	stream := c.stream
	stopRecv := c.stopRecv
	c.mu.Unlock()

	if stream == nil {
		return
	}

	go func() {
		buf := make([]byte, defaultRecvBufferSize)

		for {
			select {
			case <-stopRecv:
				return
			default:
			}

			n, err := stream.Read(buf)
			if err != nil {
				c.handleDisconnect(err, disconnected)
				return
			}

			frameBytes := buf[:n]

			typeName, msgID, err := protocol.GetFrameType(frameBytes)
			if err != nil {
				c.handleParseError(ErrFrameParseFailed)
				continue
			}

			switch typeName {
			case protocol.FrameTypeMeowOK:
				c.handleMeowOK(msgID)

			case protocol.FrameTypeError:
				var ef protocol.ErrorFrame
				if json.Unmarshal(frameBytes, &ef) != nil {
					log(LogError, "failed to parse ERROR frame")
					continue
				}
				c.handleErrorFrame(ef)

			case protocol.FrameTypeData:
				c.handleDataFrame(frameBytes)

			case protocol.FrameTypeStatusRes:
				c.handleStatusResFrame(frameBytes)

			default:
				c.handleParseError(ErrUnknownFrameType)
			}
		}
	}()
}

func (c *KittyClient) handleDisconnect(err error, disconnected chan struct{}) {
	c.mu.Lock()
	dh := c.disconnectHandler
	c.mu.Unlock()

	if dh != nil {
		dh(err)
	} else {
		log(LogError, "disconnected: %v", err)
	}

	select {
	case <-disconnected:
	default:
		close(disconnected)
	}
}

func (c *KittyClient) handleParseError(err error) {
	c.mu.Lock()
	eh := c.errHandler
	c.mu.Unlock()

	if eh != nil {
		eh("PARSE_ERROR", err.Error())
	} else {
		log(LogError, "parse error: %v", err)
	}
}
