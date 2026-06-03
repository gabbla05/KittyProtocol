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
  Lekki, bezpieczny i nowoczesny komunikator oparty o własny protokół QUIC + E2EE.
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

# 🐾 O projekcie

## ✨ Czym jest Kitty Protocol?

**Kitty Protocol** to autorski, lekki i bezpieczny protokół komunikacyjny zaprojektowany do szybkiej, prywatnej wymiany wiadomości tekstowych.

**Meowssenger** bazuje na architekturze **Client–Server**, gdzie centralny **Hub** pełni wyłącznie rolę routera:

- 🚫 nie przechowuje historii wiadomości,
- 🔐 nie posiada kluczy użytkowników,
- 🛡️ nie ma dostępu do treści komunikacji.

Całość została zaprojektowana z naciskiem na:
- prywatność,
- wydajność,
- bezpieczeństwo,
- nowoczesny transport sieciowy.

---

# Instrukcja uruchomienia

Dostępna w pliku [INSTRUCTION.md](markdowns/INSTRUCTION.md)

---

# 🖼️ Architektura projektu

<div align="center">

| Komponent | Opis |
|---|---|
| 🐈 Hub | Router QUIC/TLS 1.3 |
| 💻 CLI Client | Terminalowy klient protokołu |
| 🎨 GUI Client | Desktop app w Go + Fyne |

</div>

---

# 🐈 Hub (Server)

Centralny element infrastruktury:

- obsługa QUIC + TLS 1.3,
- zarządzanie sesjami,
- logowanie i rejestracja,
- przekazywanie zaszyfrowanych ramek,
- rate limiting,
- replay protection,
- brak dostępu do treści wiadomości.

<div align="center">
  <img src="assets/img/hub.png" width="700">
</div>

---

# 💻 CLI Client

Nowoczesny klient terminalowy:

- pełna obsługa protokołu,
- kolorowe logi,
- szybkie komendy,
- idealny do debugowania,
- niski narzut systemowy.

### Dostępne komendy

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

Graficzna aplikacja desktopowa napisana w **Go + Fyne**.

## Funkcje GUI

* ✨ nowoczesny interfejs,
* 🔐 ekran logowania,
* 💬 widok czatu,
* 🎨 custom Pink Theme,
* 📦 osadzone zasoby i fonty,
* ⚡ szybkie działanie natywne.

<div align="center">
  <img src="assets/img/meowssenger_gui_client.png" max-height: 40vw>
</div>

---

# ⚡ Najważniejsze funkcje

<div align="center">

| Funkcja               | Opis                                  |
| --------------------- | ------------------------------------- |
| 🔐 E2EE               | End-to-End Encryption                 |
| ⚡ QUIC                | Nowoczesny transport niskich opóźnień |
| 🧩 JSON Frames        | Czytelna struktura ramek              |
| 🛡️ Replay Protection | Ochrona przed Replay Attack           |
| 🔑 TOFU               | Trust On First Use                    |
| 🚀 TLS 1.3            | Pełne szyfrowanie transportu          |

</div>

---

# 🧠 Technologie

<p align="center">

<img src="https://skillicons.dev/icons?i=go,postgres,docker,linux,git,github">

</p>

---

# 📂 Struktura projektu

```bash
gui_src/
 ├── main.go
 ├── resources/
 ├── state/
 ├── theme/
 └── views/
```

---

# 📚 Dokumentacja

Pełna dokumentacja protokołu: [KittyProtocol.pdf](docs/KittyProtocol.pdf)

Zawiera:

* architekturę,
* diagramy przepływu,
* model bezpieczeństwa,
* strukturę ramek,
* opis komunikacji,
* analizę bezpieczeństwa.

---

# ☁️ Deployment — Azure Hub

## Szybkie uruchomienie

### 1️⃣ Utwórz VM

```bash
Ubuntu 22.04 (B1s)
```

### 2️⃣ Otwórz port UDP

```bash
9999/UDP
```

### 3️⃣ Skonfiguruj ENV

```bash
export KITTY_DB_DSN="postgres://kitty:kittypass@localhost:5432/kittyhub?sslmode=disable"

export KITTY_INTERCEPT_ADDR="0.0.0.0:9999"
```

### 4️⃣ Uruchom Hub

```bash
go run ./cmd/hub
```

**Pełna instrukcja konfiguracji huba dostępna [tutaj](markdowns/RemoteHubConfig.md)**

---

# 👥 Autorzy

<div align="center">

| Autor             | GitHub                                                     |
| ----------------- | ---------------------------------------------------------- |
| Gabriela Błaut    | [@gabbla05](https://github.com/gabbla05)                   |
| Michał Brzeziński | [@Michal-Brzezinski](https://github.com/Michal-Brzezinski) |
| Aleksandra Gołek  | [@styliana](https://github.com/styliana)                   |

</div>

---

# 📄 Licencja

Projekt udostępniany na **[licencji MIT](LICENSE)**.

---

<div align="center">

<h2>Zostaw gwiazdkę jeśli projekt Ci się podoba! ⭐</h2>

Made with 💖 by Kitty Protocol Team

</div>
