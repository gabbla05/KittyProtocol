![Kitty logo](kitty_logo.png "Shiprock, New Mexico by Beau Rogers")

 <img src="https://flagcdn.com/gb.svg" width="16"> - 
 English  
 <img src="https://flagcdn.com/pl.svg" width="16"> - 
[Polski](README.pl.md)

# Kitty Communication Protocol
A proprietary protocol for peer-to-peer host communication.

**Version:** 1.0.0  
**Status:** Technical Documentation (Proprietary)  
**Repository:** [github.com/gabbla05/KittyProtocol](https://github.com/gabbla05/KittyProtocol)

<br><br>
___
# 1. Purpose and Scope of the Protocol

### 1.1. Intended Use
The primary goal of **KittyProtocol** is to enable secure, ephemeral, and near-instant exchange of short text messages and real-time presence statuses. The protocol was designed with immediate data delivery in mind, eliminating the need for persistent conversation history storage on the server.

### 1.2. Problems Addressed
The protocol serves as a lightweight alternative to traditional messaging applications, focusing on two key areas:
* **Minimizing latency:** Using the QUIC transport layer drastically reduces connection establishment time.
* **Mobility:** The protocol solves the problem of session drops during network changes (e.g., switching from Wi-Fi to LTE) by leveraging the *Connection ID* mechanism.

### 1.3. Communication Model
Communication takes place in a **Peer-to-Peer (P2P)** model. Users establish connections directly with each other, ensuring the highest level of privacy.
* **Central server (Signaling Server):** Serves an auxiliary role (authentication, status management, IP address exchange).
* **Data transmission:** Occurs bypassing the server once a direct QUIC tunnel has been established.

<br><br>
___
# 2. Technical Assumptions

### 2.1. Transport Layer: QUIC

Choosing the QUIC protocol provides:
* **0-RTT Connection Establishment:** Reduced connection setup time.
* **No Head-of-Line Blocking:** Thanks to independent data streams.
* **Native encryption:** Using the TLS 1.3 standard.

### 2.2. Data Format: JSON
**JSON** was chosen as the encoding format, enabling:
* Easy debugging and developer readability.
* Elimination of byte-order (*Endianness*) issues.
* High flexibility and ease of protocol extension.

### 2.3. Reliability and Limits
| Feature | Mechanism |
| :--- | :--- |
| **Acknowledgements** | Logical ACK (`MEOW_OK`) after full JSON validation. |
| **Rate Limiting** | Max 10 messages/s per user; 100/min per IP address. |
| **Delivery ACK** | `Delivered` status sent after confirming the recipient received the message. |
| **Timeouts** | 15 seconds for authorization; 60 seconds of inactivity before session close. |


<br><br>
___
# 3. Message Structure

### 3.1. Message Types
* `HELLO` – Session initialization.
* `AUTH` – User authentication.
* `DATA` – Payload transmission (text/ASCII).
* `MEOW_OK` – Acknowledgement (ACK).
* `HISS_NAH` – Negative acknowledgement (NAK).
* `ERROR` – Error communication.
* `PING` – Session keep-alive.
* `BYE` – Session termination.

### 3.2. Example Data Frame

```json
{
  "type": "DATA",
  "msg_id": 104,
  "token": "abc-123-xyz",
  "payload": "Meow! Hello from KittyProtocol.",
  "hmac": "f2d93a8c3d7a..."
}
```

- Validation (what is considered a format error): The application rigorously validates every incoming frame. A format error (resulting in immediate message rejection and an ERROR response with code ERR_02) is defined as:
    - Parse error: The received string is not a valid JSON document (e.g., missing braces or quotation marks).
    - Missing base fields: The JSON object lacks the `type` or `msg_id` key.
    - Unknown type: The `type` key contains a value outside the allowed list (e.g., "typo_in_name").
    - Type mismatch: E.g., the `msg_id` field was sent as a String instead of an Integer.
    - Integrity failure: The hash computed by the receiver for the `payload` field does not match the value sent in the `hmac` field.

<br><br>
___

# 4. State Model and Communication Flow

- Session description (Communication flow):
The lifecycle of every session in KittyProtocol is divided into 4 strict phases:
  - Connection Establishment (Handshake): The QUIC transport layer sets up a secure connection (TLS 1.3). The first host then initiates application-level communication by sending a `HELLO` message. The second host responds with `MEOW_OK`, opening a window for login, or `HISS_NAH` requesting the `HELLO` message to be resent.
  - Authorization (`PURR_AUTH`): The sending host transmits an `AUTH` message containing credentials (username and password). The receiving host verifies the data and, upon success, replies with `MEOW_OK` along with a generated, temporary session token.
  - Data Exchange: After obtaining the token, the sender gains permission to send `DATA` messages to other users. Each `DATA` message must include a valid token and an hmac code. The receiver verifies correctness and returns `MEOW_OK` (status: Delivered) to the sender.
  - Session Closure: Occurs in two ways: gracefully (the host sends a `BYE` message and the server closes the stream) or by failure (due to timeout or protocol error).

Example sequence diagram:
The diagram below illustrates the correct message flow from connection establishment to data exchange.


- Keep-alive and reliability mechanisms:
  - Keep-alive: To prevent inactive QUIC streams from being closed by network devices (NAT/Firewall), the client is required to send a `PING` message every 30 seconds.
  - Timeouts: * Authorization timeout: A strict 15-second limit is enforced for correctly sending the `AUTH` message from the moment `HELLO` is sent. Exceeding this time results in error ERR_03 and connection termination.
  - Session timeout: If the receiver does not receive any data (neither `DATA` nor `PING` messages) from the sender within 60 seconds, the session is unilaterally closed (the host is considered offline).
  - Retransmission / Retry: At layer 4 (QUIC), lost packet retransmission is handled fully automatically and transparently for the application. At the application level: if the sender transmits a `DATA` message and does not receive `MEOW_OK` (Application-level ACK) within 5 seconds, it retries sending the same message with a new, higher `msg_id` to avoid rejection by the anti-Replay mechanism.

<br><br>
___
# 5. Security
- Confidentiality: Natively provided by the QUIC transport layer, which enforces TLS 1.3. All communication (including JSON structure and Payload) is encrypted using AES with a key established during the handshake. This completely prevents eavesdropping by third parties.
- Integrity: Each message contains an HMAC code computed from the Payload field, allowing the receiver to verify that data has not been modified in transit. Malformed messages are immediately rejected with error ERR_02.
- Authentication: The user logs in during the `AUTH` phase using a login+password pair or a unique API key transmitted in a secure stream. Identity is verified by the Hub before access to the data exchange phase is granted.
- Authorization: After successful authentication, the server generates a temporary session token bound to the host's unique Connection ID. This token must be attached to every `DATA` message so the receiver can authorize its communication with the receiving device.
- Replay Attack Protection: The system uses unique message identifiers (`msg_id`) and timestamps (`timestamp`) with a strict tolerance of 2 seconds relative to the other host's time. Each Nonce is single-use within a session, blocking the possibility of resending a captured frame.

- Threat model: KittyProtocol assumes operation in an untrusted network environment (e.g., public Wi-Fi). Below is a summary of the main attack vectors and their mitigation status:
- Eavesdropping and Man-in-the-Middle (MitM) attacks: Natively mitigated. Using QUIC enforces TLS encryption, which completely prevents third parties from reading the data exchange (including JSON structure and Payload).
- In-transit message tampering: Mitigated. The use of HMAC codes guarantees integrity. Any attempt by an attacker to modify data will result in failed verification and immediate frame rejection (error ERR_02).
- Replay attacks: Mitigated. Thanks to a strict time window (2-second tolerance relative to the other host's time), a unique `msg_id`, and a single-use Nonce, an attacker cannot resend a previously captured (even encrypted) frame.
- Resource exhaustion (DoS / Spam attacks): Partially mitigated. Application-level Rate Limiting (e.g., 10 messages per second per user / max 100 messages per minute from one IP) prevents basic spam. The Client-Server (Hub) architecture remains inherently vulnerable to large-scale DDoS attacks targeting the central node's network infrastructure.
- Endpoint compromise (Infected Host): Out of scope. KittyProtocol does not protect against malicious software running directly on the sender's or receiver's device, which could intercept the temporary session token or read a message before encryption.

<br><br>
___
# 6. Error Handling and Connection Failure
In KittyProtocol, errors are communicated asynchronously via `ERROR` frame types, enabling immediate application response.
- Error codes and their meaning:

    | Code   | Meaning    | When it occurs    | Recommended host reaction      |
    |-------|--------------|--------------------|-----------------------------|
    | ERR_01 | Protocol Violation | State sequence mismatch (e.g., sending DATA before HELLO) | Terminate current session; resend HELLO; if recurring, report error to user            |
    | ERR_02 | Format Error | Invalid JSON; missing type or msg_id; type mismatch | Display error description; correct frame and retry; terminate session on recurring errors   |
    | ERR_03 | Authorization Timeout   | Missing AUTH within authorization limit (default 15 s)  | Display timeout; retry login; if still failing, terminate session   |
    | ERR_04 | Authentication Failed              | Invalid login credentials / API key                                           | Request correct credentials; do not retry automatically without user intervention                |
    | ERR_05 | Token Invalid                      | Token missing, expired, or mismatched with Connection ID                           | Force re-authorization (AUTH); do not accept further DATA                                         |
    | ERR_06 | HMAC Mismatch                      | HMAC does not match payload                                                    | Reject message; recommend resending with correct HMAC                                            |
    | ERR_07 | Replay Detected                    | Reuse of msg_id/Nonce outside allowed window                               | Reject frame; log incident; optionally block session if attack is suspected                  |
    | ERR_08 | Rate Limit Exceeded                | Limit exceeded (e.g., >10/s or >100/min from IP)                              | Temporarily block sender; return info on block duration; host should apply backoff           |
    | ERR_09 | Delivery Failed Recipient Offline  | Recipient unavailable and delivery not possible                              | Notify sender; do not store message; suggest resending after session resumes        |
    | ERR_10 | Session Timeout Idle               | Session considered offline after prolonged inactivity                                     | Sender should attempt to resume connection; perform full handshake if desired                       |
    | ERR_11 | Resource Exhaustion                | Server resource shortage (e.g., DDoS, out of memory)                                    | Return temporary error; sender should apply backoff and retry with jitter                          |
    | ERR_12 | Internal Server Error              | Internal server error                                                          | Notify user; retry after random delay                                                       |
    | ERR_13 | Version Mismatch                   | Unsupported protocol version                                                  | Update host or negotiate compatibility; terminate session if incompatible                       |
    | ERR_14 | Unsupported Media                  | Payload of unsupported type/size                                            | Reject; inform host of allowed types/limits                                               |
    | ERR_15 | Not Authorized                     | Insufficient permissions to perform action                                                | Request different permissions; terminate send attempt                                                     |



- Behavior after syntax/protocol errors:
  - Upon detecting a format error, the receiver immediately stops processing the faulty frame and sends back an `ERROR` message with an appropriate description.
  - Serious structural violations or recurring errors result in unilateral session closure by the receiver to protect resources.
- Connection timeouts:
  - Authorization timeout: The `PURR_AUTH` process must be completed within 15 seconds of session initialization (HELLO).
  - Idle timeout: If the server detects no activity (no `DATA` or `PING`) for 60 seconds, the client is considered offline and the session is closed.
- Connection loss during session
  - By using the QUIC transport layer, KittyProtocol provides high resilience to brief transmission interruptions and network migration (e.g., Wi-Fi to LTE) without dropping the session.
  - In the event of a permanent connection loss, the protocol provides a reconnection mechanism that allows session resumption based on the last successfully processed `msg_id`.
- Duplicates and incomplete messages
  - Duplicates: The `msg_id` field, which must be unique within a session, protects against Replay attacks and accidental message duplication (e.g., by application retransmission).
  - Incompleteness: Data stream integrity is guaranteed by the QUIC layer. At the application level, any JSON frame that fails full structural validation is rejected as a format error ERR_02.
- Limits and abuse protection
  - Rate Limiting: A limit of 10 messages per second per user is enforced to prevent spamming and overloading the routing node.
  - IP Protection: The server/hosts block traffic exceeding 100 messages per minute from a single IP address, providing a basic barrier against DoS attacks.
  - Message size: The protocol enforces a maximum size for a single JSON frame (e.g., 64 KB) to optimize parsing performance and prevent attacks based on sending unnaturally large data payloads.

- In KittyProtocol, errors are communicated via `ERROR` frames containing the fields:

  - `type: "ERROR"`
  - `Msg_id`
  - `code` – error code / short description of the error type
  - `desc` – error description
  
<br><br>
___
# 7. Scenarios (Frame Examples)

   - Scenario 1: Successful session and message exchange
   This is the standard flow from connection establishment to sending a short text message.
       1. Connection Establishment (Handshake):
       - Sender: `{"type": "HELLO", "msg_id": 1, …}`
       - Receiver: `{"type": "MEOW_OK", "msg_id": 1, …}`

       2. Authorization (PURR_AUTH):
       - Sender: 
      ` {"type": "AUTH", "msg_id": 2, "login": "user_name", "pass": "secret_purr"}`
       - Receiver: 
      ` {"type": "MEOW_OK", "msg_id": 2, "token": "abc-123-xyz", "expires": 3600}`

       3. Data Exchange (DATA):
       - Sender: 
       `{"type": "DATA", "msg_id": 3, "token": "abc-123-xyz", "to": "mordunia2", "payload": "Meow! Check out this ASCII cat:", "hmac": "f2d9..."}`
       - Receiver (to sender): 
       `{"type": "MEOW_OK", "msg_id": 3, "status": "Delivered"}`

  - Scenario 2: Format validation error (ERR_02)
      A scenario in which the sender sends a message that does not comply with the specification (e.g., incorrect length or missing required fields).
      - Sender: 
      `{"type": "DATA", "msg_id": 10, "token": "abc-123-xyz", "payload": "Text too long...", "length": 5}` (Declared length of 5 does not match actual length)
      - Receiver: 
      `{"type": "ERROR", "code": "ERR_02", "msg_id": 10, "desc": "Invalid data format or HMAC"}`

   - Scenario 3: Authorization process timeout (ERR_03)
        A scenario illustrating the protocol's strict time limits.
        - Sender: `{"type": "HELLO", "msg_id": 50}`
        - Receiver: `{"type": "MEOW_OK", "msg_id": 50}`
        (No host activity for more than 15 seconds)
        - Receiver: `{"type": "ERROR", "code": "ERR_03", "desc": "PURR_AUTH timeout reached"}`
        - Connection: Receiver unilaterally closes the QUIC session.