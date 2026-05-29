[🇬🇧 English version](README.md) | 🇵🇱 *Wersja polska*

# Kitty Protocol x Meowssenger

<div style="display: flex; gap: 5%; align-items: center; width: 100%;">
  <img src="assets/img/kitty_logo.png" style="width: 60%; height: auto;">
  <img src="gui_src/resources/assets/images/logo_z_napisem.png" style="width: 35%; height: auto;">
</div>



## 🐾 Opis projektu

**Kitty Protocol** to autorski, lekki i bezpieczny protokół komunikacyjny zaprojektowany do szybkiej, prywatnej wymiany krótkich wiadomości tekstowych.  
**Meowssenger** to system opierający się na architekturze **Client–Server**, gdzie centralny serwer (*Hub*) pełni wyłącznie rolę routera — **nie przechowuje historii**, **nie zna kluczy**, **nie ma dostępu do treści wiadomości**.

Projekt składa się z trzech głównych komponentów:
- **Hub (serwer)** — router QUIC/TLS 1.3 z obsługą sesji, autoryzacji i ochrony przed atakami.
   <img src="assets/img/hub.png" width="400">
   <br><br>

- **Klient CLI** — terminalowy interfejs użytkownika z pełną obsługą protokołu.
   <img src="assets/img/meowssenger_cli_client.png" width="400">
   <br><br>

- **Klient GUI (Fyne)** — graficzna aplikacja desktopowa zbudowana w Go + Fyne.
   <img src="assets/img/meowssenger_gui_client.png" width="400">
   <br><br>

---

## ✨ Najważniejsze funkcje

### 🔐 End‑to‑End Encryption (E2EE)
Wszystkie wiadomości są szyfrowane **wyłącznie po stronie klienta**.  
Hub nie posiada kluczy, nie może odszyfrować treści i nie przechowuje historii.

### ⚡ Transport QUIC
Protokół działa na warstwie **QUIC**:
- niskie opóźnienia,
- odporność na zmiany sieci (Wi‑Fi ↔ LTE),
- wsparcie dla 0‑RTT,
- natywne TLS 1.3.

### 🧩 Ramki JSON
Wszystkie komunikaty protokołu są przesyłane w formacie **JSON**, co ułatwia debugowanie i rozwój.

### 🛡️ Odporność na ataki
Wbudowane mechanizmy:
- ochrona przed **Replay Attack**,
- **TOFU** (Trust On First Use) dla kluczy serwera,
- ograniczanie sesji i rate‑limiting,
- pełne szyfrowanie transportowe (TLS 1.3).

---

## 🖥️ Komponenty projektu

### 🐈 Hub (serwer)
Centralny router protokołu:
- nasłuchuje na QUIC/TLS,
- zarządza sesjami,
- obsługuje logowanie i rejestrację,
- przekazuje zaszyfrowane ramki między klientami,
- nie przechowuje treści wiadomości.

### 💻 Klient CLI
Terminalowy interfejs użytkownika:
- pełna obsługa protokołu,
- komendy `/chat`, `/msg`, `/status`, `/secret`, `/logout`,
- kolorowe logi i czytelny interfejs,
- idealny do testów i debugowania.

### 🎨 Klient GUI (Fyne)
Graficzna aplikacja desktopowa:
- widok logowania,
- widok menu,
- widok czatu,
- własny motyw kolorystyczny (Pink Theme),
- zasoby (fonty, grafiki) osadzone w binarce.

Struktura GUI:
```
gui_src/
 ├── main.go
 ├── resources/
 ├── state/
 ├── theme/
 └── views/
```

---

## 📚 Dokumentacja

Pełna dokumentacja protokołu znajduje się tutaj:  
**[docs/KittyProtocol.pdf](docs/KittyProtocol.pdf)**

Zawiera:
- opis ramek,
- diagramy przepływu,
- model bezpieczeństwa,
- scenariusze komunikacji,
- strukturę transportu QUIC.

---

## 🚀 Uruchamianie Huba na Azure (skrót)

1. Utwórz VM (Ubuntu 22.04, B1s).  
2. Otwórz port **UDP 9999** w Azure Networking.  
3. Zainstaluj Go i PostgreSQL.  
4. Skonfiguruj zmienne środowiskowe:
   ```
   export KITTY_DB_DSN="postgres://kitty:kittypass@localhost:5432/kittyhub?sslmode=disable"
   export KITTY_INTERCEPT_ADDR="0.0.0.0:9999"
   ```
5. Uruchom:
   ```
   go run ./cmd/hub
   ```

Pełny samouczek wdrożenia dostępny w pliku [RemoteHubConfig.md](markdowns/RemoteHubConfig.md).

---

## 👥 Autorzy

Projekt został stworzony przez:

- [Gabriela Błaut](https://github.com/gabbla05)  
- [Michał Brzeziński](https://github.com/Michal-Brzezinski)  
- [Aleksandra Gołek](https://github.com/styliana)

---

## 📄 Licencja

Projekt udostępniany na licencji MIT.

