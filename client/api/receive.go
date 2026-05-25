package api

import (
	"encoding/json"

	"github.com/gabbla05/KittyProtocol/internal/cryptoee"
	"github.com/gabbla05/KittyProtocol/protocol"
)

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

func (c *KittyClient) StartReceiverLoop(disconnected chan struct{}) {
	c.mu.Lock()
	stream := c.stream
	replay := c.replay
	ackMgr := c.ackMgr
	stopRecv := c.stopRecv

	helloCh := c.helloCh
	authCh := c.authCh
	registerCh := c.registerCh

	chatReqCh := c.chatReqCh
	chatAcceptCh := c.chatAcceptCh
	chatRefuseCh := c.chatRefuseCh
	chatEndCh := c.chatEndCh
	chatMsgCh := c.chatMsgCh
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
				// dynamiczny disconnectHandler
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
				return
			}

			frameBytes := buf[:n]

			typeName, msgID, err := protocol.GetFrameType(frameBytes)
			if err != nil {
				c.mu.Lock()
				eh := c.errHandler
				c.mu.Unlock()

				if eh != nil {
					eh("PARSE_ERROR", err.Error())
				} else {
					log(LogError, "parse error: %v", err)
				}
				continue
			}

			switch typeName {

			// ============================================================
			//  MEOW_OK — odpowiedź na HELLO / AUTH / REGISTER
			// ============================================================
			case protocol.FrameTypeMeowOK:
				c.mu.Lock()
				currentState := c.state
				c.mu.Unlock()

				switch currentState {

				case StateHandshaking:
					helloCh <- OpResult{OK: true}
					c.setState(StateAuthenticating)

				case StateAuthenticating:
					authCh <- OpResult{OK: true}
					c.setState(StateSelectingTarget)

				case StateRegistering:
					registerCh <- OpResult{OK: true}
					c.setState(StateAuthenticating)

				default:
					if ackMgr != nil {
						ackMgr.NotifyDelivered(msgID)
					}
				}

			// ============================================================
			//  ERROR — odpowiedź na HELLO / AUTH / REGISTER
			// ============================================================
			case protocol.FrameTypeError:
				var ef protocol.ErrorFrame
				if json.Unmarshal(frameBytes, &ef) != nil {
					log(LogError, "failed to parse ERROR frame")
					continue
				}

				c.mu.Lock()
				currentState := c.state
				eh := c.errHandler
				c.mu.Unlock()

				switch currentState {

				case StateHandshaking:
					helloCh <- OpResult{OK: false, Code: ef.Code, Desc: ef.Desc}

				case StateAuthenticating:
					authCh <- OpResult{OK: false, Code: ef.Code, Desc: ef.Desc}

				case StateRegistering:
					registerCh <- OpResult{OK: false, Code: ef.Code, Desc: ef.Desc}

				default:
					if eh != nil {
						eh(ef.Code, ef.Desc)
					} else {
						log(LogError, "server error %s: %s", ef.Code, ef.Desc)
					}
				}

			// ============================================================
			//  DATA — E2EE chat payload (lub inne dane aplikacyjne)
			// ============================================================
			case protocol.FrameTypeData:
				var df protocol.DataFrame
				if json.Unmarshal(frameBytes, &df) != nil {
					log(LogError, "failed to parse DATA frame")
					continue
				}

				if replay != nil && replay.MarkAndCheck(df.MsgID) {
					continue
				}

				c.mu.Lock()
				kEnc, kMac, ok := c.getKeysForPeer(df.Sender)
				appHandler := c.appHandler
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

				// Spróbuj zinterpretować jako ChatFrame
				var probe chatFrameProbe
				if err := json.Unmarshal([]byte(plaintext), &probe); err == nil && probe.Type != "" {
					switch probe.Type {

					case "CHAT_REQUEST":
						if chatReqCh != nil {
							chatReqCh <- ChatRequestEvent{From: probe.From}
						}

					case "CHAT_ACCEPT":
						if chatAcceptCh != nil {
							chatAcceptCh <- ChatAcceptEvent{From: probe.From}
						}

					case "CHAT_REFUSE":
						var p chatRefusePayload
						_ = json.Unmarshal(probe.Payload, &p)
						if chatRefuseCh != nil {
							chatRefuseCh <- ChatRefuseEvent{From: probe.From, Reason: p.Reason}
						}

					case "CHAT_END":
						var p chatEndPayload
						_ = json.Unmarshal(probe.Payload, &p)
						if chatEndCh != nil {
							chatEndCh <- ChatEndEvent{From: probe.From, Reason: p.Reason}
						}

					case "TEXT_MESSAGE":
						var p textMessagePayload
						_ = json.Unmarshal(probe.Payload, &p)
						if chatMsgCh != nil {
							chatMsgCh <- ChatMessageEvent{From: probe.From, Text: p.Text}
						}

					default:
						log(LogWarn, "unknown chat frame type: %s", probe.Type)
					}

					// Chat‑frame obsłużony
					continue
				}

				// Fallback: inne dane aplikacyjne
				if appHandler != nil {
					appHandler(df.Sender, []byte(plaintext))
				}

			// ============================================================
			//  STATUS_RES
			// ============================================================
			case protocol.FrameTypeStatusRes:
				var sf protocol.StatusResFrame
				if json.Unmarshal(frameBytes, &sf) != nil {
					log(LogError, "failed to parse STATUS_RES")
					continue
				}

				c.mu.Lock()
				sh := c.statusHandler
				c.mu.Unlock()

				if sh != nil {
					sh(sf.Target, sf.Status)
				} else {
					log(LogInfo, "status: %s is %s", sf.Target, sf.Status)
				}
			}
		}
	}()
}
