 # KittyProtocol – Developer Architecture Overview

 ## 1. High‑Level Architecture
 Projekt składa się z trzech głównych warstw:
 - **Protocol Layer** – definicje ramek, formaty JSON, reguły protokołu.
 - **Client Layer** – implementacja klienta QUIC + E2EE + TOFU + ACK.
 - **Hub Layer** – serwer QUIC, routing wiadomości, autoryzacja, sesje.

 Każda warstwa jest odseparowana i ma jasno określone odpowiedzialności.

 ---

 ## 2. Protocol Layer (`/protocol`)
 Zawiera wyłącznie definicje ramek i funkcje ich parsowania.

 ### Najważniejsze struktury:
 - `BaseFrame{Type, MsgID}`
 - `HelloFrame`, `AuthFrame`, `DataFrame`, `StatusResFrame`, `ErrorFrame`, `MeowOkFrame`

 ### Odpowiedzialność:
 - Format JSON ramek.
 - Walidacja pól.
 - Funkcje `ParseXFrame(raw)` i `GetFrameType(raw)`.

 ### Co tu NIE występuje:
 - Żadna logika sieciowa.
 - Żadne szyfrowanie.
 - Żadne operacje na sesjach.

 ---

 ## 3. Client Layer (`/client`)
 Warstwa klienta składa się z trzech podwarstw:
 - **API** – logika protokołu i QUIC.
 - **APP** – logika aplikacji (menu, chat).
 - **UI** – interfejs użytkownika (CLI).

 ### 3.1. `client/api`
 Odpowiada za:
 - QUIC Connect / Stream.
 - TLS 1.3 + TOFU (Trust On First Use).
 - E2EE (kEnc, kMac) – ustawiane przez `SetSharedSecret()`.
 - Wysyłanie ramek: `SendAuth`, `SendMessage`, `SendGetStatus`, `SendBye`.
 - Odbiór ramek: `StartReceiverLoop()`.
 - Ping loop: `StartPingLoop()`.
 - ACK subsystem: `AckManager`.

 #### Najważniejsze moduły:
 - `tls.go` – hardened TLS 1.3 config + VerifyConnection hook.
 - `tofu.go` – pinning certyfikatu serwera.
 - `e2ee.go` – KDF i ustawianie kluczy szyfrujących.
 - `ack.go` – zarządzanie ACK/timeoutami.
 - `send.go` – wysyłanie ramek.
 - `receive.go` – odbiór ramek i dispatch.

 #### Przykładowy przepływ:
 1. `Connect()` → QUIC handshake.
 2. `WaitForHelloOK()` → potwierdzenie od Huba.
 3. `SendAuth()` → wysłanie loginu/hasła.
 4. `WaitForAuthOK()` → potwierdzenie logowania.
 5. `SendMessage()` → wysłanie DATA.
 6. `AckManager` → czeka na MEOW_OK lub timeout.

 ---

 ## 3.2. `client/app`
 Warstwa logiki aplikacji – **nie** drukuje nic sama, tylko korzysta z UI.

 ### Pliki:
 - `app.go` – konstruktor + interfejs UI.
 - `menu.go` – główne menu aplikacji.
 - `chat.go` – logika sesji czatu.

 ### Odpowiedzialność:
 - Interpretacja komend użytkownika.
 - Wywoływanie metod API.
 - Obsługa rozłączenia.

 ### Komendy:
 - `/status <user>` – wysyła GET_STATUS.
 - `/chat <user>` – wchodzi w tryb czatu.
 - `/logout` – wylogowuje się, kończy sesję i zamyka program klienta.

 ---

 ## 3.3. `client/ui_cli`
 Warstwa wejścia/wyjścia – **tylko** terminal.

 ### Odpowiedzialność:
 - `ReadLine()`, `Println()`, `Printf()`.
 - Pobieranie loginu/hasła.
 - Pobieranie sekretu E2EE.
 - Implementacja `AckEventHandler` (drukuje Delivered/Timeout).

 ### Co tu NIE występuje:
 - Żadna logika protokołu.
 - Żadna logika aplikacji.

 ---

 ## 4. Hub Layer (`/hub`)
 Serwer QUIC odpowiedzialny za routing i zarządzanie sesjami.

 ### 4.1. `main.go`
 - Startuje listener QUIC.
 - Ładuje certyfikaty (certmanager).
 - Obsługuje sygnały OS.
 - Dla każdej sesji wywołuje `handleClient()`.

 ---

 ## 4.2. Handlery (`handler_*.go`)

 ### `handler_dispatcher.go`
 - Odbiera raw JSON.
 - Wyciąga typ ramki (`GetFrameType`).
 - Wywołuje odpowiedni handler:
   - `handleHello`
   - `handleAuth`
   - `handlePing`
   - `handleData`
   - `handleGetStatus`
   - `handleBye`

 ### `handler_auth.go`
 - Parsuje AUTH.
 - Zatrzymuje timer AUTH.
 - Weryfikuje login/hasło.
 - Tworzy sesję (`SessionManager.Add`).
 - Wysyła MEOW_OK("Logged in").

 ### `handler_data.go`
 - Parsuje DATA.
 - Sprawdza sesję.
 - Rate limiting.
 - Replay protection.
 - Aktualizacja aktywności.
 - Routing do odbiorcy (`routeData`).
 - Wysyła MEOW_OK ACK.

 ### `handler_status.go`
 - Parsuje GET_STATUS.
 - Sprawdza czy użytkownik online.
 - Wysyła STATUS_RES.

 ### `handler_bye.go`
 - Usuwa sesję.
 - Zamyka połączenie.

 ### `handler_hello.go`
 - Wysyła MEOW_OK("Ready for auth").
 - Startuje timer AUTH.

 ### `handler_ping.go`
 - Aktualizuje LastActive.

 ---

 ## 4.3. Routing (`router.go`)
 - Znajduje sesję odbiorcy.
 - Forwarduje DATA.
 - Obsługuje błędy (ERR_10, ERR_15).

 ---

 ## 4.4. Errors (`errors.go`)
 - Funkcja `sendError(stream, code, desc)`.
 - Wysyła standardową ramkę ERROR.

 ---

 ## 4.5. Session Manager (`internal/protection`)
 - Przechowuje aktywne sesje.
 - Rate limiting.
 - Replay protection.
 - Idle timeout (LastActive).

 ---

 ## 5. Certmanager (`internal/certmanager`)
 - Generuje poprawne self‑signed certyfikaty ECDSA.
 - Obsługuje SAN (DNS + IP).
 - Zgodny z TLS 1.3 i TOFU.

 ---

 ## 6. Przykładowy pełny przepływ (Client → Hub → Client)

 ### 1. Klient łączy się z Hubem (QUIC + TLS 1.3 + TOFU)
 - Klient weryfikuje certyfikat serwera.
 - Jeśli pierwszy raz → zapisuje certyfikat (TOFU).

 ### 2. HELLO
 - Klient wysyła HELLO.
 - Hub odpowiada MEOW_OK("Ready for auth").

 ### 3. AUTH
 - Klient wysyła login/hasło.
 - Hub weryfikuje i tworzy sesję.
 - Hub wysyła MEOW_OK("Logged in").

 ### 4. Chat
 - Klient ustawia sekret E2EE.
 - Klient wysyła DATA.
 - Hub forwarduje DATA do odbiorcy.
 - Odbiorca wysyła MEOW_OK (ACK).
 - AckManager klienta kończy timer.

 ### 5. GET_STATUS
 - Klient wysyła GET_STATUS.
 - Hub odpowiada STATUS_RES("online"/"offline").

 ### 6. BYE
 - Klient wysyła BYE.
 - Hub usuwa sesję.

 ---

 ## 7. Najważniejsze zasady architektury
 - **UI nie zawiera logiki** – tylko I/O.
 - **APP nie zna protokołu** – tylko workflow.
 - **API nie zna UI** – tylko QUIC + protokół.
 - **Hub nie zna klienta** – tylko ramki JSON.
 - **Protocol nie zna niczego** – tylko struktury.

 ---

 ## 8. Co jest częścią protokołu?
 - Format ramek JSON.
 - Typy ramek (HELLO, AUTH, DATA, GET_STATUS, STATUS_RES, MEOW_OK, ERROR, BYE).
 - Reguły: AUTH przed DATA, ACK po DATA, timeouty.
 - Routing DATA.
 - Kody błędów (ERR_XX).

 ---

 ## 9. Co NIE jest częścią protokołu?
 - UI (CLI).
 - Menu i komendy.
 - Format logowania.
 - Implementacja certmanagera.
 - E2EE (to warstwa aplikacji, nie protokołu).

 ---

 ## 10. Podsumowanie
 KittyProtocol to:
 - QUIC + TLS 1.3 + TOFU
 - JSON‑based messaging protocol
 - E2EE po stronie klienta
 - Routing + sesje po stronie Huba
 - Czysty podział warstw: Protocol → API → App → UI
