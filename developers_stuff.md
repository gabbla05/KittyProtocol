 # KittyProtocol – Developer Documentation (Stage 1)
 
 KittyProtocol is a lightweight, secure, QUIC‑based communication protocol designed for
 ephemeral, end‑to‑end encrypted messaging. This document summarizes the current state
 of the implementation, architecture, module responsibilities, and development workflow.
 It reflects the system exactly as it exists after completion of Gołson’s tasks and before
 Message Broker (MB) and E2EE (Stage 2) are implemented.
 
 ---
 
 ## 1. Architecture Overview
 
 KittyProtocol follows a Client–Server (Hub) model:
 
 - Clients establish QUIC/TLS 1.3 connections to the Hub.
 - Hub acts as a router (Message Broker) and identity manager.
 - All payloads are encrypted end‑to‑end (E2EE) — Hub cannot read message contents.
 - Hub stores no message history and no cryptographic keys.
 - QUIC ensures low latency, mobility, and automatic retransmission.
 
 The protocol defines strict phases:
 
 1. HELLO → server responds MEOW_OK("Ready for auth")
 2. AUTH → server validates credentials
 3. GET_STATUS (future)
 4. DATA exchange (E2EE)
 5. MEOW_OK acknowledgements
 
 ---
 
 ## 2. Repository Structure
 
 ```
 ├── certs/                  # TLS certificates for QUIC
 ├── client/                 # Client application (CLI)
 │   ├── main.go             # Entry point
 │   ├── hello_auth.go       # HELLO/AUTH logic
 │   ├── input.go            # User input helpers
 │   ├── ping.go             # Keep‑alive PING sender
 │   ├── ack.go              # ACK tracking and timers
 ├── hub/                    # Hub (server)
 │   ├── main.go             # Listener + QUIC setup
 │   ├── handler.go          # Frame processing loop
 │   ├── auth_flow.go        # HELLO/AUTH logic
 │   ├── errors.go           # Standardized error frames
 ├── internal/
 │   ├── auth/               # Mock DB + bcrypt verification
 │   ├── clientutils/        # TruncateMessage, ACK timers
 │   └── protection/         # RateLimiter, SessionManager, AuthTimer
 ├── protocol/
 │   └── frames.go           # UniversalFrame + JSON parsing
 └── README.md
 ```
 
 ---
 
 ## 3. Completed Tasks (Gołson)
 
 ### 3.1 Task 6 – Authentication (Mock DB)
 
 - Implemented bcrypt‑based credential verification.
 - Users: `alice` and `bob` with hashed passwords.
 - AUTH flow fully functional:
   - HELLO → MEOW_OK("Ready for auth")
   - AUTH(user, pass) → MEOW_OK("Logged in") or ERR_04
 
 ### 3.2 Task 9 – Hub Protection Mechanisms
 
 - **Auth Timeout (20s)**  
   If AUTH is not received within 20 seconds after HELLO, Hub sends ERR_03 and closes the connection.
 
 - **Idle Timeout (60s)**  
   SessionManager removes inactive sessions and closes QUIC connections.
 
 - **Rate Limiting (Token Bucket)**  
   - 10 messages per second per user.
   - Prevents spam and DoS.
 
 ### 3.3 Task 14 – Client Edge Logic
 
 - **TruncateMessage** ensures plaintext stays below safe limit before encryption.
 - **ACK Timer (5s)** marks messages as undelivered if MEOW_OK is not received.
 
 ---
 
 ## 4. Current Behavior (Stage 1)
 
 ### 4.1 Successful AUTH
 
 ```
 [Server]: Ready for auth
 Login: alice
 Hasło: secret
 [Server]: Logged in
 ```
 
 ### 4.2 Sending Messages
 
 ```
 > Hello Bob
 [Delivered] msg_id=1776478265666
 ```
 
 ### 4.3 Idle Timeout
 
 ```
 [Protection] Idle Timeout: alice. Removing session.
 ```
 
 Client sees:
 
 ```
 [Client] Connection closed: timeout: no recent network activity
 ```
 
 This is correct — QUIC reports the connection closure triggered by Hub.
 
 ### 4.4 Missing Routing (Expected)
 
 If routing is not implemented yet (Task 7), ACK timer expires:
 
 ```
 [Status] Wiadomość 1776478...: Niedostarczono (timeout)
 ```
 
 This is correct until MB implements message forwarding.
 
 ---
 
 ## 5. Test Coverage (Stage 1)
 
 Tests implemented:
 
 - `RateLimiter` (token refill, exhaustion)
 - `SessionManager` (idle cleanup)
 - `ParseFrame` (valid/invalid JSON)
 - `TruncateMessage` (boundary conditions)
 
 Current coverage:
 
 ```
 internal/clientutils     50%
 internal/protection      63%
 protocol                 75%
 ```
 
 Additional tests planned:
 
 - SessionManager concurrency tests
 - ACK timer cancellation tests
 - HELLO/AUTH integration tests (after routing)
 
 ---
 
 ## 6. Next Steps (Stage 2 – MB + Gaba)
 
 ### 6.1 Task 7 – Message Broker (MB)
 
 - Implement routing:
   - Hub receives DATA from Alice
   - Looks up Bob’s session
   - Forwards DATA to Bob
   - Bob sends MEOW_OK
   - Hub forwards MEOW_OK to Alice
 
 ### 6.2 Task 11 – E2EE (MB)
 
 - AEAD_Encrypt(K_enc, nonce=msg_id)
 - HMAC(K_mac, cipher|msg_id|target)
 - AEAD_Decrypt on receiver side
 
 ### 6.3 Task 12 – Client State Machine (Gaba)
 
 - DISCONNECTED → HELLO → AUTH → ESTABLISHED
 - Handle ERR_03, ERR_04, ERR_09 transitions
 
 ### 6.4 Task 10 – Frame Models (Gaba)
 
 - Replace UniversalFrame with typed frames:
   - HelloFrame
   - AuthFrame
   - DataFrame
   - ErrorFrame
   - StatusFrame
 
 ---
 
 ## 7. Running the System
 
 ### 7.1 Start Hub
 ```
 go run hub/main.go
 ```
 
 ### 7.2 Start Client
 ```
 go run client/main.go
 ```
 
 ### 7.3 Multiple Clients
 ```
 go run client/main.go -port 9001
 go run client/main.go -port 9002
 ```
 
 ---
 
 ## 8. Summary
 
 Stage 1 of KittyProtocol is complete:
 
 - Authentication works (bcrypt)
 - QUIC/TLS 1.3 stable
 - Session management stable
 - Rate limiting stable
 - ACK timers stable
 - Test suite implemented
 
 System is ready for Stage 2:
 
 - Message Broker (routing)
 - E2EE
 - State machine
 - Typed frames
 
 This document will be extended in Stage 2.
