# Działanie protokołu pomiędzy sieciami

## Cel dokumentu
Celem dokumentu jest opisanie, jak KittyProtocol został przetestowany w środowisku wielosieciowym
(klient za NAT ↔ Hub na Azure ↔ drugi klient), jak wygląda komunikacja QUIC/TLS na poziomie sieci,
oraz potwierdzenie, że protokół działa poprawnie, stabilnie i bezpiecznie.

## Środowisko testowe

### Hub (Azure VM)
• System: Ubuntu na Azure B1s  
• Interfejs: `eth0`  
• Adres prywatny: `10.0.0.4/24`  
• Nasłuch: `0.0.0.0:9999/UDP`  
• Proces: `go run ./hub`  

### Klient 1 (Alice – laptop domowy)
• Sieć: domowy internet, NAT/CGNAT operatora  
• Proces: `go run ./client`  
• Łączy się z publicznym IP Azure VM  

### Klient 2 (Bob – druga instancja)
• Również `go run ./client`  
• Łączy się z tym samym Hubem  

### Protokół transportowy
• QUIC (quic-go)  
• TLS 1.3  
• Port: `9999/UDP`  

## Scenariusz testowy

### 1. Uruchomienie Huba
Hub startuje na Azure i nasłuchuje na `0.0.0.0:9999`.

### 2. Klient Alice (laptop)
• Wysyła `HELLO`  
• Otrzymuje `MEOW_OK`  
• Wysyła `AUTH` (`alice` / `secret`)  
• Otrzymuje `MEOW_OK` („Logged in”)  

### 3. Klient Bob (drugi klient)
• Analogiczny handshake  
• Obie sesje widoczne w `SessionManager`  

### 4. Wymiana wiadomości DATA
• Bob → Alice: `DATA`  
• Alice → Bob: `DATA`  
• Hub przekazuje zaszyfrowane payloady  
• Nadawca dostaje `MEOW_OK` jako ACK  

### 5. Zapytania GET_STATUS
• Alice: `/status bob` → `GET_STATUS`  
• Hub: sprawdza `globalSessions`  
• Hub: odsyła `STATUS_RES` (`online` / `offline`)  
• Po wyjściu Boba → `offline`  

### 6. Zamykanie sesji
• Klient wysyła `BYE`  
• Hub usuwa sesję  
• Idle Timeout usuwa nieaktywne sesje po 60s  

## Analiza ruchu sieciowego (tshark / QUIC)

### 1. Przechwytywanie ruchu
Komenda:
`sudo tshark -i eth0 -f "udp port 9999" -w /tmp/kitty-cross-network.pcapng`

Podgląd:
`sudo tshark -r /tmp/kitty-cross-network.pcapng`

### 2. Co widać w zrzucie

#### Handshake QUIC
Pierwsze pakiety mają typ:
`QUIC Initial, CRYPTO`

To oznacza handshake TLS 1.3 w ramach QUIC.

#### Przejście do szyfrowania
Po handshake wszystkie pakiety mają formę:
`QUIC Protected Payload (KP0)`

To oznacza:
• Cała treść KittyProtocol (HELLO, AUTH, DATA, GET_STATUS, STATUS_RES, PING, BYE)  
• jest zaszyfrowana i niewidoczna w tsharku.  

#### Dwa handshaki
W zrzucie widać dwa bloki handshake:
• sesja Alice ↔ Hub  
• sesja Bob ↔ Hub  

#### Ruch w obie strony
Widać pakiety:
• `185.188.117.206 → 10.0.0.4` (klient → Hub)  
• `10.0.0.4 → 185.188.117.206` (Hub → klient)  

To potwierdza poprawne działanie przez NAT/CGNAT.

#### Duże pakiety (1408–1483 bajty)
To zaszyfrowane:
• DATA  
• STATUS_RES  
• ACK  
• PING  

QUIC pakuje dane do MTU ~1500 bajtów.

### 3. Wnioski bezpieczeństwa
• Payload jest całkowicie niewidoczny (TLS 1.3 + QUIC).  
• Podsłuchujący widzi tylko metadane (IP, port, rozmiar).  
• E2EE dodatkowo szyfruje treść wiadomości — nawet Hub jej nie zna.  

## Zachowanie protokołu w środowisku wielosieciowym

### NAT / CGNAT po stronie klienta
• Klient nie ma publicznego IP  
• Nie może przyjmować połączeń  
• Ale może inicjować QUIC → Hub  

### Publiczny Hub na Azure
• Osiągalny z dowolnej sieci  
• Centralny router ramek DATA  
• Odpowiada na GET_STATUS  

### Stabilność
• QUIC utrzymuje sesję mimo NAT  
• PING utrzymuje aktywność  
• Idle Timeout czyści sesje  

## Podsumowanie
• KittyProtocol działa poprawnie między różnymi sieciami (Azure ↔ laptop).  
• QUIC/TLS 1.3 działa stabilnie i bezpiecznie.  
• Payload jest zaszyfrowany i niewidoczny w tsharku.  
• Routing DATA działa poprawnie.  
• GET_STATUS działa poprawnie.  
• BYE i Idle Timeout działają zgodnie z założeniami.  
• Protokół jest gotowy do dalszego rozwoju (GUI, backend użytkowników, testy obciążeniowe).  
