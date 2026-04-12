hub:  go run hub/main.go
clients:  go run client/main.go


<img width="608" height="250" alt="image" src="https://github.com/user-attachments/assets/b19229c1-ca5a-47ce-b1e6-e49c5519362f" />


Aktualne zachowanie aplikacji
(Status: Oczekiwanie na ACK
- Opis: Po wysłaniu wiadomości przez klienta, w konsoli pojawia się komunikat o oczekiwaniu na potwierdzenie, który po 5 sekundach zgłasza timeout.Przyczyna techniczna: Jest to poprawne działanie timera z Taska 14. Ponieważ Hub nie posiada jeszcze zaimplementowanego routingu (Task 7), wiadomość nie jest przekazywana dalej, a co za tym idzie – nadawca nigdy nie otrzymuje ramki MEOW_OK.
- - Wymagane wdrożenia: Task 7 (MB): Uruchomienie przekazywania strumieni między sesjami w Hubie. Task 12 (Gaba): Implementacja maszyny stanów obsługującej logiczne potwierdzenia.)

# 🐈 KittyProtocol - Dokumentacja Deweloperska (Etap 1)

Ten dokument zawiera podsumowanie aktualnej architektury, struktury plików oraz zrealizowanych funkcjonalności w ramach projektu KittyProtocol.

## 1. Cel i Model Komunikacji
KittyProtocol to bezpieczny, bezstanowy system wymiany wiadomości oparty na **QUIC** i **TLS 1.3**. 
* **Model**: Architektura Klient-Serwer (C/S), gdzie centralny **Hub** pełni rolę routera (Message Broker)
* **E2EE**: Treść wiadomości jest szyfrowana end-to-end; Hub widzi jedynie metadane niezbędne do routingu.

## 2. Struktura Projektu (Rafinacja Gołsona)
Projekt został podzielony na moduły wewnątrz folderu `internal/`, aby odizolować logikę transportową od zabezpieczeń.

```bash
├── certs/              # Certyfikaty TLS 1.3 dla QUIC
├── client/             # Implementacja aplikacji klienckiej
├── hub/                # Główny serwer (Router) 
├── internal/           # Logika biznesowa (Taski Gołsona) 
│   ├── auth/           # Weryfikacja poświadczeń (Mock DB + bcrypt) 
│   ├── clientutils/    # Narzędzia pomocnicze klienta (Truncation, Timery) 
│   └── protection/     # Mechanizmy ochronne Huba (Rate Limiter, Session Manager) 
└── protocol/           # Definicje ramek JSON (Single Source of Truth) 

```
3. Status Implementacji (Taski Gołsona)
W ramach infrastruktury i ochrony zrealizowano następujące zadania:
- Task 6: Autoryzacja i Mock DB:Stworzono moduł internal/auth/auth.go obsługujący weryfikację użytkowników (Alice/Bob) przy użyciu bcrypt.
- Task 9: Ochrona Huba:
--Auth Timeout (20s): Hub zamyka połączenie, jeśli autoryzacja nie nastąpi w wyznaczonym czasie. -- Idle Timeout (60s): SessionManager automatycznie usuwa nieaktywne sesje z pamięci RAM. -- Rate Limiting: Algorytm Token Bucket ograniczający ruch do 10 wiadomości/sekundę na użytkownika.
- Task 14: Logika Klienta: -- Implementacja TruncateMessage – automatyczne przycinanie tekstu do 2048 bajtów.
- Mechanizm ACK Timer (5s) – oczekiwanie na potwierdzenie MEOW_OK przed zmianą statusu wiadomości.

4. Instrukcja Uruchomienia
- Zależności: Zainstaluj biblioteki: go mod tidy.
- Certyfikaty: Upewnij się, że pliki .pem znajdują się w folderze certs/.
- Start Huba: go run hub/main.go.
- Start Klienta: go run client/main.go (użyj flagi -port dla wielu instancji).
