package api

import (
	"encoding/json"

	"github.com/gabbla05/KittyProtocol/internal/cryptoee"
	"github.com/gabbla05/KittyProtocol/protocol"
)

// StartReceiverLoop launches a background goroutine responsible for reading
// all incoming frames from the QUIC stream. This is the only reader for the
// stream; all other components must communicate through higher‑level APIs.
func (c *KittyClient) StartReceiverLoop(disconnected chan struct{}) {
	c.mu.Lock()
	stream := c.stream
	replay := c.replay
	ackMgr := c.ackMgr
	stopRecv := c.stopRecv
	errHandler := c.errHandler
	statusHandler := c.statusHandler
	disconnectHandler := c.disconnectHandler
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
				if disconnectHandler != nil {
					disconnectHandler(err)
				} else {
					log(LogError, "disconnected: %v", err)
				}

				select {
				case <-disconnected:
				default:
					close(disconnected)
				}
				return
			}

			typeName, msgID, err := protocol.GetFrameType(buf[:n])
			if err != nil {
				if errHandler != nil {
					errHandler("PARSE_ERROR", err.Error())
				} else {
					log(LogError, "parse error: %v", err)
				}
				continue
			}

			switch typeName {

			case protocol.FrameTypeMeowOK:
				if ackMgr != nil {
					ackMgr.NotifyDelivered(msgID)
				}

			case protocol.FrameTypeError:
				var ef protocol.ErrorFrame
				if json.Unmarshal(buf[:n], &ef) == nil {
					if errHandler != nil {
						errHandler(ef.Code, ef.Desc)
					} else {
						log(LogError, "server error %s: %s", ef.Code, ef.Desc)
					}
				} else {
					log(LogError, "failed to parse ERROR frame")
				}

			case protocol.FrameTypeData:
				var df protocol.DataFrame
				if json.Unmarshal(buf[:n], &df) != nil {
					log(LogError, "failed to parse DATA frame")
					continue
				}

				if replay != nil && replay.MarkAndCheck(df.MsgID) {
					continue
				}

				c.mu.Lock()
				handler := c.appHandler
				kEnc, kMac, ok := c.getKeysForPeer(df.Sender)
				c.mu.Unlock()

				if !ok {
					log(LogWarn, "no shared secret for %s", df.Sender)
					continue
				}

				plaintext, err := cryptoee.DecryptAndVerifyWithKeys(
					df.MsgID, df.Target, df.Payload, df.MAC, kEnc, kMac,
				)
				if err != nil {
					log(LogError, "E2EE error: %v", err)
					continue
				}

				if handler != nil {
					handler(df.Sender, []byte(plaintext))
				}

			case protocol.FrameTypeStatusRes:
				var sf protocol.StatusResFrame
				if json.Unmarshal(buf[:n], &sf) != nil {
					log(LogError, "failed to parse STATUS_RES")
					continue
				}

				if statusHandler != nil {
					statusHandler(sf.Target, sf.Status)
				} else {
					log(LogInfo, "status: %s is %s", sf.Target, sf.Status)
				}
			}
		}
	}()
}
