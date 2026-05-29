🇬🇧 *English version* | [🇵🇱 Wersja polska](README.pl.md) 
# Kitty Protocol x Meowssenger

<div style="display: flex; gap: 5%; align-items: center; width: 100%;">
  <img src="assets/img/kitty_logo.png" style="width: 60%; height: auto;">
  <img src="gui_src/resources/assets/images/logo_z_napisem.png" style="width: 35%; height: auto;">
</div>


## 🐾 Project Overview

**Kitty Protocol** is a custom, lightweight and secure communication protocol designed for fast, private exchange of short text messages.  
**Meowssenger** is the system that follows a **Client–Server** architecture, where the central server (*Hub*) acts purely as a router — it **stores no history**, **knows no keys**, and **cannot access message content**.

The project consists of three main components:
- **Hub (server)** — QUIC/TLS 1.3 router with session management, authentication and attack protection.
  
    <img src="assets/img/hub.png" width="400">
  <br><br>

- **CLI Client** — terminal-based user interface with full protocol support.
  <img src="assets/img/meowssenger_cli_client.png" width="400">
<br><br>

- **GUI Client (Fyne)** — graphical desktop application built with Go + Fyne.
  <img src="assets/img/meowssenger_gui_client.png" width="400">

---

## ✨ Key Features

### 🔐 End‑to‑End Encryption (E2EE)
All messages are encrypted **exclusively on the client side**.  
The Hub has no keys, cannot decrypt content, and stores no message history.

### ⚡ QUIC Transport
The protocol is built on **QUIC**, providing:
- low latency,
- resilience to network changes (Wi‑Fi ↔ LTE),
- 0‑RTT support,
- native TLS 1.3 security.

### 🧩 JSON Frames
All protocol messages are transmitted in **JSON format**, simplifying debugging and development.

### 🛡️ Attack Resistance
Built‑in mechanisms include:
- protection against **Replay Attacks**,
- **TOFU** (Trust On First Use) for server key pinning,
- session limiting and rate‑limiting,
- full TLS 1.3 transport encryption.

---

## 🖥️ Project Components

### 🐈 Hub (Server)
The central router of the protocol:
- listens on QUIC/TLS,
- manages sessions,
- handles login and registration,
- forwards encrypted frames between clients,
- stores no message content.

### 💻 CLI Client
Terminal-based user interface:
- full protocol support,
- commands `/chat`, `/msg`, `/status`, `/secret`, `/logout`,
- readable logs and structured output,
- ideal for testing and debugging.

### 🎨 GUI Client (Fyne)
Graphical desktop application:
- authentication view,
- menu view,
- chat view,
- custom color theme (Pink Theme),
- embedded assets (fonts, images).

GUI structure:
```
gui_src/
 ├── main.go
 ├── resources/
 ├── state/
 ├── theme/
 └── views/
```

---

## 📚 Documentation

Full protocol documentation is available here:  
**[docs/KittyProtocol-EN.pdf](docs/KittyProtocol-EN.pdf)**

It includes:
- frame definitions,
- flow diagrams,
- security model,
- communication scenarios,
- QUIC transport structure.

---

## 🚀 Running the Hub on Azure (short version)

1. Create a VM (Ubuntu 22.04, B1s).  
2. Open **UDP 9999** in Azure Networking.  
3. Install Go and PostgreSQL.  
4. Configure environment variables:
   ```
   export KITTY_DB_DSN="postgres://kitty:kittypass@localhost:5432/kittyhub?sslmode=disable"
   export KITTY_INTERCEPT_ADDR="0.0.0.0:9999"
   ```
5. Run:
   ```
   go run ./cmd/hub
   ```

A full deployment guide is available in the markdown [RemoteHubConfig.md](markdowns/RemoteHubConfig.md) file.

---

## 👥 Authors

- [Gabriela Błaut](https://github.com/gabbla05)  
- [Michał Brzeziński](https://github.com/Michal-Brzezinski)  
- [Aleksandra Gołek](https://github.com/styliana)

---

## 📄 License

This project is released under the MIT License.  
See the **LICENSE** file for details.