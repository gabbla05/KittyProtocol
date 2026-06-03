<div align="center">

<h1 align="center" style="
    border-bottom: none;
    font-family: 'Poppins', 'Segoe UI', sans-serif;
    font-size: 3rem;
    font-weight: 800;
    background: linear-gradient(90deg, #ff4fa3, #ff85c1, #ffb6e6);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    text-shadow: 0 0 12px rgba(255, 105, 180, 0.35);
    letter-spacing: 1px;
    margin-bottom: 0;
">
    ✨ Kitty Protocol x Meowssenger ✨
</h1>

<p align="center">
  <img src="assets/img/kitty_logo.png" height="120" alt="Kitty Logo">
  &nbsp;&nbsp;&nbsp;&nbsp;
  <img src="gui_src/resources/assets/images/logo_z_napisem.png" height="120" alt="Kitty Protocol">
</p>

<p align="center">
  <strong>Cute • Secure • Decentralized Messaging 🐾</strong>
</p>

<p align="center">
  A lightweight, secure and modern messenger powered by a custom QUIC + E2EE protocol.
</p>

---

<p align="center">

<img src="https://img.shields.io/badge/Protocol-Kitty%20Protocol-ff69b4?style=for-the-badge">
<img src="https://img.shields.io/badge/Transport-QUIC-e75480?style=for-the-badge">
<img src="https://img.shields.io/badge/Encryption-E2EE-ff1493?style=for-the-badge">
<img src="https://img.shields.io/badge/Language-Go-f06292?style=for-the-badge&logo=go">
<img src="https://img.shields.io/badge/GUI-Fyne-ff85c1?style=for-the-badge">

</p>

<p align="center">

[![en](https://img.shields.io/badge/lang-en-blue.svg)](README.md)
[![pl](https://img.shields.io/badge/lang-pl-red.svg)](README.pl.md)
<img src="https://img.shields.io/badge/license-MIT-ff69b4.svg">
<img src="https://img.shields.io/github/stars/gabbla05/KittyProtocol?style=social&dummy=1">

</p>

</div>

---

# About the Project

## What is Kitty Protocol?

**Kitty Protocol** is a custom, lightweight and secure communication protocol designed for fast, private exchange of text messages.

**Meowssenger** follows a **Client–Server** architecture where the central **Hub** acts purely as a router:

- 🚫 stores no message history,
- 🔐 has no access to user keys,
- 🛡️ cannot read message contents.

The project was designed with a strong focus on:

- privacy,
- performance,
- security,
- modern networking technologies.

---

# Execution instruction

Available in file: [INSTRUCTION.md](markdowns/INSTRUCTION.md)

---

# Project Architecture

<div align="center">

| Component | Description |
|---|---|
| 🐈 Hub | QUIC/TLS 1.3 Router |
| 💻 CLI Client | Terminal-based protocol client |
| 🎨 GUI Client | Desktop application built with Go + Fyne |

</div>

---

# Hub (Server)

The central component of the infrastructure:

- QUIC + TLS 1.3 transport,
- session management,
- authentication and registration,
- encrypted frame routing,
- rate limiting,
- replay protection,
- zero access to message content.

<div align="center">
  <img src="assets/img/hub.png" width="700">
</div>

---

# CLI Client

Modern terminal-based client:

- full protocol support,
- colorful logs,
- fast command system,
- ideal for debugging,
- lightweight and efficient.

### Available Commands

```bash
/status <user>
/secret <user>
/chat <user> 
/accept <user>
/refuse <user>
/msg <text>
/end
/logout
/menu
/help

````

<div align="center">
  <img src="assets/img/meowssenger_cli_client.png" width="750">
</div>

---

# 🎨 GUI Client (Fyne)

Desktop application built with **Go + Fyne**.

## GUI Features

* ✨ modern user interface,
* 🔐 authentication screen,
* 💬 chat view,
* 🎨 custom Pink Theme,
* 📦 embedded fonts and assets,
* ⚡ native performance.

<div align="center">
  <img src="assets/img/meowssenger_gui_client.png">
</div>

---

# Key Features

<div align="center">

| Feature               | Description                       |
| --------------------- | --------------------------------- |
| 🔐 E2EE               | End-to-End Encryption             |
| ⚡ QUIC                | Low-latency modern transport      |
| 🧩 JSON Frames        | Human-readable frame structure    |
| 🛡️ Replay Protection | Protection against replay attacks |
| 🔑 TOFU               | Trust On First Use                |
| 🚀 TLS 1.3            | Secure transport encryption       |

</div>

---

# Technologies

<p align="center">

<img src="https://skillicons.dev/icons?i=go,postgres,docker,linux,git,github">

</p>

---

# 📂 Project Structure

```bash
gui_src/
 ├── main.go
 ├── resources/
 ├── state/
 ├── theme/
 └── views/
```

---

# 📚 Documentation

Full protocol documentation:

**[KittyProtocol-EN.pdf](docs/KittyProtocol-EN.pdf)**

Includes:

* protocol architecture,
* flow diagrams,
* security model,
* frame definitions,
* communication scenarios,
* QUIC transport design.

---

# ☁️ Azure Hub Deployment

## Quick Start

### 1️⃣ Create a VM

```bash
Ubuntu 22.04 (B1s)
```

### 2️⃣ Open UDP Port

```bash
9999/UDP
```

### 3️⃣ Configure Environment Variables

```bash
export KITTY_DB_DSN="postgres://kitty:kittypass@localhost:5432/kittyhub?sslmode=disable"

export KITTY_INTERCEPT_ADDR="0.0.0.0:9999"
```

### 4️⃣ Run the Hub

```bash
go run ./cmd/hub
```

**Full deployment guide available [here](markdowns/RemoteHubConfig.md)**

---

# 👥 Authors

<div align="center">

| Author            | GitHub                                                     |
| ----------------- | ---------------------------------------------------------- |
| Gabriela Błaut    | [@gabbla05](https://github.com/gabbla05)                   |
| Michał Brzeziński | [@Michal-Brzezinski](https://github.com/Michal-Brzezinski) |
| Aleksandra Gołek  | [@styliana](https://github.com/styliana)                   |

</div>

---

# 📄 License

This project is distributed under the **[MIT License](LICENSE)**.

---

<div align="center">

<h2>Leave a star if you like the project! ⭐</h2>

Made with 💖 by Kitty Protocol Team

</div>
