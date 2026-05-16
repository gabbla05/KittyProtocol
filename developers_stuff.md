# KAŻDY MUSI SOBIE ZROBIĆ PLIK .ENV I DO NIEGO SKOPIOWAĆ ZAWARTOŚĆ .ENV.EXAMPLE LUB PODMIENIĆ WEFŁUG PREFERENCJI

# KittyProtocol – Developer Documentation (Stage 1 → Stage 2 Update)

KittyProtocol is a lightweight, secure, QUIC‑based communication protocol designed for
ephemeral, end‑to‑end encrypted messaging. This document summarizes the current state
of the implementation, architecture, module responsibilities, and development workflow.

This version reflects the system after completing Stage 1 and Stage 2:
- QUIC/TLS 1.3 transport
- Session management and rate limiting
- Strict JSON validation
- Typed frames
- Message Broker (routing)
- End‑to‑End Encryption (E2EE)
- Replay‑protection (client + hub)
- ACK timers
- Keep‑alive (PING)
- Logging

The structure of the original Stage 1 document has been preserved.
Only sections requiring updates have been modified.

---

## 1. Architecture Overview

KittyProtocol follows a Client–Server (Hub) model:

- Clients establish QUIC/TLS 1.3 connections to the Hub.
- Hub acts as a router (Message Broker) and identity manager.
- All payloads are encrypted end‑to‑end (E2EE) — Hub cannot read message contents.
- Hub stores no message history and no cryptographic keys.
- QUIC ensures low latency, mobility, and automatic retransmission.

Protocol phases:

1. HELLO → server responds MEOW_OK("Ready for auth")
2. AUTH → server validates credentials
3. GET_STATUS (in progress)
4. DATA exchange (E2EE)
5. MEOW_OK acknowledgements

Stage 2 additions:
- Full E2EE (AEAD + HMAC)
- Replay‑protection (ERR_06 + silent drop)
- Typed frames
- Strict JSON parser
- Keep‑alive (PING)
- Logging

---

## 2. Repository Structure

```
.
├── certs/
│   ├── cert.pem
│   └── key.pem
├── client/
│   ├── ack.go
│   ├── bye.go
│   ├── hello_auth.go
│   ├── input.go
│   ├── main.go
│   ├── ping.go
│   ├── replay.go
│   ├── replay_test.go
│   └── state.go
├── cmd/
│   ├── hash/
│   │   └── main.go
│   └── mock_server/
│       └── main.go
├── docs/
│   ├── protocol_schema.json
│   └── types.ts
├── documentation/
│   ├── KittyProtocol-EN.pdf
│   └── KittyProtocol.pdf
├── hub/
│   ├── auth_flow.go
│   ├── errors.go
│   ├── handler.go
│   ├── main.go
│   └── router.go
├── internal/
│   ├── auth/
│   │   └── auth.go
│   ├── certmanager/
│   │   ├── certmanager.go
│   │   └── certmanager_test.go
│   ├── clientutils/
│   │   ├── ack_timer.go
│   │   ├── truncate.go
│   │   └── truncate_test.go
│   ├── cryptoee/
│   │   ├── cryptoee_test.go
│   │   ├── decrypt.go
│   │   ├── encrypt.go
│   │   └── keys.go
│   └── protection/
│       ├── limiter.go
│       ├── limiter_test.go
│       ├── replay.go
│       ├── replay_test.go
│       ├── session.go
│       ├── session_manager.go
│       └── session_manager_test.go
├── protocol/
│   ├── frames.go
│   └── frames_test.go
├── README.md
└── README.pl.md
```

---

## 3. Completed Tasks (Stage 1 → Stage 2)

### 3.1 Authentication (Mock DB)
- bcrypt‑based credential verification
- HELLO → AUTH → MEOW_OK flow

### 3.2 QUIC/TLS 1.3 + Session Management
- TLS 1.3 enforced (certmanager)
- ALPN: kitty-quic-v1
- Safe stream handling
- SessionManager with idle timeout

### 3.3 Protection Mechanisms
- Auth Timeout (20s)
- Idle Timeout (60s)
- Token Bucket rate limiting (10 msg/s)

### 3.4 Strict JSON Validation
- `GetFrameType()` validates JSON structure and required fields
- Missing or malformed fields → ERR_02

### 3.5 Typed Frames
- HelloFrame, AuthFrame, DataFrame, ErrorFrame, StatusFrame

### 3.6 Message Broker (Routing)
- Forwarding DATA from sender → receiver
- Forwarding MEOW_OK back to sender
- Intelligent routing based on `target`

### 3.7 End‑to‑End Encryption (E2EE)
- AEAD encryption using msg_id as nonce
- HMAC integrity protection
- Hub cannot decrypt or forge messages

### 3.8 Replay Protection
- Hub: ERR_06 for duplicate msg_id
- Client: silent drop of duplicate DATA
- TTL‑based eviction + memory cap

### 3.9 Client State Machine
- DISCONNECTED → HANDSHAKING → AUTH → ESTABLISHED

### 3.10 ACK Timers
- 5‑second timeout for MEOW_OK

### 3.11 Keep‑alive (PING)
- Client sends PING every 30s
- Hub responds MEOW_OK

### 3.12 Logging
- Diagnostic logs for all frames (Hub + Client)

---

## 4. Current Behavior (Stage 2)

### 4.1 Successful AUTH
```
[Server]: Ready for auth
Login: alice
Password: secret
[Server]: Logged in
```

### 4.2 Sending Messages (E2EE)
```
> Hello Bob
[Delivered] msg_id=1776478265666
```

### 4.3 Replay Handling
- Client → Hub replay → ERR_06
- Hub → Client replay → silently ignored

### 4.4 Idle Timeout
Hub:
```
[Protection] Idle Timeout: alice. Removing session.
```

Client:
```
[Client] Connection closed: timeout: no recent network activity
```

---

## 5. Test Coverage (Stage 2)

Tests implemented:
- RateLimiter
- SessionManager
- Strict JSON parser
- CertManager
- ReplayDetector (client + hub)
- cryptoee (encrypt/decrypt roundtrip)
- ACK timers

Coverage:
```
internal/clientutils     ~50%
internal/protection      ~80%
internal/certmanager     100%
protocol                 100%
cryptoee                 100%
```

---

## 6. Next Steps (Stage 3 – Remaining Tasks)

### 6.1 KIT‑12 — Headless test script
- Automated end‑to‑end test: Hub + 2 clients

### 6.2 KIT‑14 — Load test (max clients)
- Determine maximum number of concurrent clients

### 6.3 KIT‑15 — CPU performance test
- Analyze routing performance vs CPU cores

### 6.4 KIT‑16 — Negative tests
- ERR_04 (bad password)
- ERR_15/ERR_16 (offline user)

### 6.5 KIT‑20 — GET_STATUS / STATUS_RES
- Online/offline presence system

### 6.6 KIT‑17 — PostgreSQL AUTH backend
- Replace mock DB with real database

### 6.7 KIT‑33 — Permission system
- Allow‑list for incoming messages (anti‑spam)

### 6.8 GUI tasks (Meowssenger)
- KIT‑22: Main window
- KIT‑23: Search Username (GET_STATUS)
- KIT‑24: Local E2EE module
- KIT‑25: Local conversation state
- KIT‑26: Send button → DATA frame
- KIT‑28: Happy Path tests
- KIT‑29: Security tests

---

## 7. Running the System

### 7.1 Start Hub
```
go run ./hub
```

### 7.2 Start Client
```
go run ./client
```

### 7.3 Multiple Clients
```
go run ./client -port 9001
go run ./client -port 9002
```

---

## 8. Summary

Stage 1 → Stage 2 transition is complete:
- QUIC/TLS 1.3 stable
- Session management stable
- Rate limiting stable
- Strict JSON validation
- Typed frames
- Message Broker routing
- E2EE implemented
- Replay‑protection implemented
- ACK timers
- Keep‑alive
- Logging

Remaining work (Stage 3):
- Automated tests
- Load/performance tests
- GET_STATUS
- PostgreSQL AUTH
- Permission system
- GUI (Meowssenger)

This document will be extended in Stage 3.

# STAGE 3

## 1. Overview
This document describes all changes introduced in the KittyProtocol project
from the moment GET_STATUS / STATUS_RES frames were implemented up to the full
rebuild of the security layer, TLS, configuration, and the certificate trust
mechanism (TOFU).

These changes include:
- user status routing,
- improvements to client and hub architecture,
- introduction of environment variables,
- elimination of hard‑coded IP,
- full QUIC/TLS encryption,
- introduction of TOFU (Trust On First Use),
- cleanup of developer tools.

## 2. GET_STATUS and STATUS_RES
### 2.1. Purpose
The GET_STATUS mechanism allows the client to ask the Hub whether a given user
is currently logged in and has an active session.

### 2.2. Implementation
- The client sends a GET_STATUS frame with the `target` field.
- The Hub checks in SessionManager whether an active session exists.
- The Hub returns STATUS_RES with `status: "online"` or `"offline"`.

### 2.3. Behavior
- After a user logs out (BYE or Idle Timeout), the status becomes offline.
- The mechanism works fully independently from the E2EE layer.

## 3. Elimination of hard‑coded IP and introduction of .env
### 3.1. Problem
Previously, the Hub address was hard‑coded in the client.
This prevented:
- environment switching,
- deployment on different machines,
- local and remote testing,
- automation.

### 3.2. Solution
An environment variable was introduced:
`KITTY_HUB_ADDR`

If not set, the client uses:
`127.0.0.1:9999`

### 3.3. .env loader
The client automatically loads the `.env` file at startup.

## 4. TLS – from InsecureSkipVerify to full security
### 4.1. Problem
Early versions of the client used:
`InsecureSkipVerify: true`

which meant:
- no server certificate verification,
- vulnerability to MITM,
- no guarantee the client connects to the real Hub.

### 4.2. Goal
Provide:
- full encryption,
- server authentication,
- MITM resistance,
- no need to manually copy certificates.

## 5. TOFU – Trust On First Use
### 5.1. Why TOFU?
Certificates generated by certmanager are self‑signed.
There is no CA, so classical TLS verification is not possible.

TOFU works like SSH:
- first connection → client saves the server certificate,
- subsequent connections → client compares the certificate with the saved one,
- if the certificate changes → MITM warning and connection aborted.

### 5.2. Implementation
After `DialAddr()` the client retrieves the certificate:
`connState := conn.ConnectionState()`

If `trusted_cert.pem` does not exist:
- the client saves the server certificate.

If it exists:
- the client compares the DER with the server certificate,
- if different → connection blocked.

### 5.3. Advantages
- full automation,
- no manual certificate copying,
- MITM resistance,
- ideal for self‑signed certificates.

## 6. Changes in certmanager
### 6.1. SAN / DNSNames
Certmanager generates certificates with:
`DNSNames: ["kitty-hub"]`

This enables:
- correct TLS behavior,
- compatibility with ServerName in the client,
- unambiguous server identification.

## 7. Project structure after changes
### 7.1. Client
- `client/main.go` – state logic, TOFU, QUIC
- `client/tls.go` – TLS configuration for TOFU
- `client/*.go` – frames, E2EE, replay protection

### 7.2. Hub
- `hub/main.go` – QUIC listener
- `hub/handler.go` – frame routing
- `hub/auth_flow.go` – AUTH
- `hub/router.go` – DATA, GET_STATUS

### 7.3. Internal
- `internal/certmanager` – certificate generation
- `internal/protection` – limiter, replay, session manager
- `internal/cryptoee` – E2EE

### 7.4. Tools
- `cmd/mock_server` – dev tool
- `cmd/hash` – hash generator (dev only)

## 8. Security tests
### 8.1. Test 1 – first connection
- the client saves the server certificate,
- creates `trusted_cert.pem`.

### 8.2. Test 2 – subsequent connection
- the client compares certificates,
- if matching → OK.

### 8.3. Test 3 – certificate change (MITM)
- the client detects the change,
- blocks the connection.

### 8.4. Test 4 – two independent clients
- run in two directories,
- each has its own pinned certificate.

## 9. Summary
We introduced:
- GET_STATUS / STATUS_RES,
- full status routing,
- environment variables,
- .env loader,
- elimination of hard‑coded IP,
- full QUIC/TLS encryption,
- TOFU (Trust On First Use),
- automatic certificate pinning,
- MITM detection,
- cleanup of developer tools.

KittyProtocol is now:
- secure,
- MITM‑resistant,
- configurable,
- deployment‑ready,
- aligned with security best practices.
