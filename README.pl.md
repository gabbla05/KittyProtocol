![Kitty logo](kitty_logo.png "Shiprock, New Mexico by Beau Rogers")

 <img src="https://flagcdn.com/pl.svg" width="16"> - 
 Polish  
 <img src="https://flagcdn.com/gb.svg" width="16"> - 
[English](README.md)


# Protokół komunikacyjny Kitty 
Autorski protokół do komunikacji pomiędzy hostami w architekturze peer2peer.

**Wersja:** 1.0.0  
**Status:** Dokumentacja Techniczna (Proprietary)  
**Repozytorium:** [github.com/gabbla05/KittyProtocol](https://github.com/gabbla05/KittyProtocol)




<br><br>
___
# 1. Cel i zakres protokołu

### 1.1. Przeznaczenie
Głównym celem **KittyProtocol** jest umożliwienie bezpiecznej, bezstanowej (*ephemeral*) i błyskawicznej wymiany krótkich wiadomości tekstowych oraz statusów obecności w czasie rzeczywistym. Protokół został zaprojektowany z myślą o natychmiastowym dostarczaniu danych, eliminując potrzebę trwałego przechowywania historii konwersacji na serwerze.

### 1.2. Rozwiązywane problemy
Protokół stanowi lekką alternatywę dla tradycyjnych komunikatorów, skupiając się na dwóch kluczowych obszarach:
* **Minimalizacja opóźnień:** Wykorzystanie warstwy transportowej QUIC drastycznie redukuje czas nawiązywania połączeń.
* **Mobilność:** Protokół rozwiązuje problem zrywania sesji przy zmianach sieci (np. przejście z Wi-Fi na LTE) dzięki wykorzystaniu mechanizmu *Connection ID*.

### 1.3. Model komunikacji
Komunikacja odbywa się w modelu **Peer-to-Peer (P2P)**. Użytkownicy nawiązują połączenia bezpośrednio między sobą, co gwarantuje najwyższy poziom prywatności. 
* **Centralny serwer (Signaling Server):** Pełni rolę pomocniczą (uwierzytelnianie, zarządzanie statusami, wymiana adresów IP).
* **Transmisja danych:** Odbywa się z pominięciem serwera po zestawieniu bezpośredniego tunelu QUIC.

<br><br>
___
# 2. Założenia techniczne

### 2.1. Warstwa transportowa: QUIC

Wybór protokołu QUIC zapewnia:
* **0-RTT Connection Establishment:** Skrócenie czasu nawiązywania połączenia.
* **Brak blokowania linii (Head-of-Line Blocking):** Dzięki niezależnym strumieniom danych.
* **Natywne szyfrowanie:** Wykorzystanie standardu TLS 1.3.

### 2.2. Format danych: JSON
Jako format kodowania wybrano **JSON**, co pozwala na:
* Łatwe debugowanie i czytelność dla programisty.
* Eliminację problemów z kolejnością bajtów (*Endianness*).
* Wysoką elastyczność i łatwość rozszerzania protokołu.

### 2.3. Niezawodność i limity
| Cecha | Mechanizm |
| :--- | :--- |
| **Potwierdzenia** | Logiczne ACK (`MEOW_OK`) po pełnej walidacji JSON. |
| **Rate Limiting** | Max 10 wiadomości/s na użytkownika; 100/min na adres IP. |
| **Delivery ACK** | Status `Delivered` wysyłany po upewnieniu się, że odbiorca odebrał komunikat. |
| **Timeouty** | 15 sekund na autoryzację; 60 sekund bezczynności do zamknięcia sesji. |


<br><br>
___
# 3. Struktura komunikatów

### 3.1. Typy komunikatów
* `HELLO` – Inicjalizacja sesji.
* `AUTH` – Uwierzytelnianie użytkownika.
* `DATA` – Przesyłanie ładunku (tekst/ASCII).
* `MEOW_OK` – Potwierdzenie (ACK).
* `HISS_NAH` – Negatywne potwierdzenie (NAK).
* `ERROR` – Komunikacja błędów.
* `PING` – Utrzymanie sesji (Keep-alive).
* `BYE` – Zakończenie sesji.

### 3.2. Przykład ramki danych

```json
{
  "type": "DATA",
  "msg_id": 104,
  "token": "abc-123-xyz",
  "payload": "Meow! Hello from KittyProtocol.",
  "hmac": "f2d93a8c3d7a..."
}
```

- Walidacja (co uznajemy za błąd formatu): Aplikacja rygorystycznie sprawdza każdą przychodzącą ramkę. Za błąd formatu (skutkujący natychmiastowym odrzuceniem wiadomości i odesłaniem komunikatu ERROR z kodem ERR_02) uznaje się:
    - Błąd parsowania: Otrzymany ciąg znaków nie jest poprawnym dokumentem JSON (np. brakujące klamry, cudzysłowy).
    - Brak pól bazowych: W obiekcie JSON brakuje klucza type lub msg_id.
    - Nieznany typ: Klucz type zawiera wartość spoza dozwolonej listy (np. "typo_w_nazwie").
    - Niezgodność typów zmiennych: Np. pole msg_id zostało przesłane jako ciąg znaków (String) zamiast liczby całkowitej (Integer).
    - Brak integralności: Wyliczony przez odbiorcę skrót dla pola payload nie zgadza się z wartością przesłaną w polu hmac.

<br><br>
___

# 4. Model stanów i przebieg komunikacji

- Opis sesji (Przebieg komunikacji):
Cykl życia każdej sesji w protokole KittyProtocol dzieli się na 4 rygorystyczne fazy:
  - Nawiązanie połączenia (Handshake): Warstwa transportowa QUIC zestawia bezpieczne połączenie (TLS 1.3). Następnie pierwszy host inicjuje komunikację aplikacyjną, wysyłając komunikat `HELLO`. Drugi host odpowiada potwierdzeniem `MEOW_OK`, otwierając okno na logowanie lub HISS_NAH prosząc o ponowne przesłanie komunikatu `HELLO`.
  - Autoryzacja (`PURR_AUTH`): Host nadawcy przesyła komunikat AUTH zawierający poświadczenia (login i hasło). Host odbiorcy weryfikuje dane i w przypadku sukcesu odsyła `MEOW_OK` z wygenerowanym, tymczasowym tokenem sesyjnym.
  - Wymiana danych: Po uzyskaniu tokenu, nadawca zyskuje uprawnienia do wysyłania komunikatów `DATA` do innych użytkowników. Każda wiadomość `DATA` musi zawierać ważny token oraz kod hmac. Odbiorca weryfikuje poprawność i zwraca nadawcy `MEOW_OK` (status: Delivered) .
  - Zamknięcie sesji: Odbywa się na dwa sposoby: świadomie (host wysyła komunikat `BYE`, a serwer zamyka strumień) lub awaryjnie (przez przekroczenie timeoutu lub błąd protokołu).

Przykładowy diagram sekwencji:
Przedstawiony schemat obrazuje prawidłowy przepływ wiadomości od nawiązania połączenia po wymianę danych.


- Mechanizmy podtrzymania i niezawodności:
  - Keep-alive: Aby zapobiec zamykaniu nieaktywnych strumieni QUIC przez urządzenia sieciowe (NAT/Firewall), wprowadzono zasadę wysyłania komunikatu `PING` przez klienta co 30 sekund.
  - Timeouty: * Timeout autoryzacji: Wprowadzono rygorystyczny limit 15 sekund na poprawne przesłanie komunikatu `AUTH` od momentu wysłania `HELLO`. Przekroczenie tego czasu skutkuje błędem ERR_03 i przerwaniem połączenia .
  - Timeout sesji: Jeśli odbiorca nie otrzyma od nadawcy żadnych danych (komunikatów `DATA` lub `PING`) w ciągu 60 sekund, następuje jednostronne zamknięcie sesji (uznanie hosta za offline).
  - Retransmisja / Retry: Na poziomie warstwy 4 (QUIC) retransmisja zagubionych pakietów realizowana jest w pełni automatycznie i przezroczyście dla aplikacji. Na poziomie aplikacyjnym: jeśli nadawca wyśle komunikat `DATA` i w ciągu 5 sekund nie otrzyma od odbiorcy zwrotnego `MEOW_OK` (Application-level ACK), ponawia próbę wysłania tej samej wiadomości z nowym, wyższym msg_id, aby uniknąć odrzucenia przez mechanizm anty-Replay.

<br><br>
___
# 5. Bezpieczeństwo
- Poufność: Realizowana natywnie przez warstwę transportową QUIC, która wymusza stosowanie protokołu TLS 1.3. Całość komunikacji (w tym struktura JSON i Payload), jest szyfrowana algorytmem AES z kluczem ustalonym podczas handshake’u . Całkowicie uniemożliwia to podsłuch przez osoby trzecie.
- Integralność: Każdy komunikat zawiera kod HMAC wyliczany z pola Payload, co pozwala odbiorcy na weryfikację, czy dane nie zostały zmodyfikowane w locie. Błędne komunikaty są natychmiast odrzucane z błędem ERR_02.
- Uwierzytelnienie: Użytkownik loguje się w fazie `AUTH` przy użyciu pary login+hasło lub unikalnego klucza API przesyłanego w bezpiecznym strumieniu. Tożsamość jest weryfikowana przez Hub przed dopuszczeniem do fazy wymiany danych.
- Autoryzacja: Po poprawnym uwierzytelnieniu, serwer generuje tymczasowy token sesyjny powiązany z unikalnym Connection ID hosta. Token ten musi być dołączany do każdego komunikatu `DATA`, aby odbiorca autoryzował jego komunikację z urządzeniem odbierającym.
- Ochrona przed Replay Attack: System wykorzystuje unikalne identyfikatory wiadomości (`msg_id`) oraz znaczniki czasu (`timestamp`) z rygorystyczną tolerancją wynoszącą 2 sekundy względem czasu drugiego hosta. Każdy Nonce jest jednorazowy w obrębie sesji, co blokuje możliwość ponownego przesłania przechwyconej ramki.

- Model zagrożeń: Protokół KittyProtocol zakłada działanie w niezaufanym środowisku sieciowym (np. publiczne sieci Wi-Fi). Poniżej zestawiono główne wektory ataków oraz status ich obsługi:
- Podsłuch i ataki Man-in-the-Middle (MitM): Mitygowane natywnie. Wykorzystanie QUIC wymusza szyfrowanie TLS, co całkowicie uniemożliwia odczytanie wymiany danych (w tym struktury JSON i Payloadu) przez osoby trzecie.
- Modyfikacja wiadomości w locie (Tampering): Mitygowane. Zastosowanie kodów HMAC gwarantuje integralność. Każda próba modyfikacji danych przez atakującego zakończy się błędną weryfikacją i natychmiastowym odrzuceniem ramki (błąd ERR_02).
- Ataki Replay: Mitygowane. Dzięki zastosowaniu rygorystycznego okna czasowego (tolerancja 2 sekundy względem czasu drugiego hosta/innych hostów), unikalnego msg_id oraz jednorazowego Nonce'a, atakujący nie może ponownie wysłać raz przechwyconej (nawet zaszyfrowanej) ramki.
- Wyczerpanie zasobów (Ataki DoS / Spam): Częściowo mitygowane. Aplikacyjny Rate Limiting (np. 10 komunikatów na sekundę na użytkownika / max 100 wiadomości na minutę z jednego IP) zapobiega prostemu spamowaniu. Architektura Client-Server (Hub) pozostaje jednak naturalnie podatna na masowe ataki DDoS na infrastrukturę sieciową centralnego węzła.
- Przejęcie punktu końcowego (Zainfekowany Host): Poza zakresem protokołu. KittyProtocol nie chroni przed złośliwym oprogramowaniem działającym bezpośrednio na urządzeniu nadawcy lub odbiorcy, które mogłoby przechwycić tymczasowy token sesyjny lub odczytać wiadomość przed jej zaszyfrowaniem.

<br><br>
___
# 6. Obsługa błędów i awarii połączenia
W protokole KittyProtocol błędy są komunikowane asynchronicznie za pomocą ramek typu `ERROR`, co pozwala na natychmiastową reakcję aplikacji.
- Kody błędów i ich znaczenie:

    | Kod   | Znaczenie    | Kiedy występuje    | Zalecana reakcja hosta      |
    |-------|--------------|--------------------|-----------------------------|
    | ERR_01 | Protocol Violation | Niezgodność z sekwencją stanów (np. wysłanie DATA przed HELLO) | Przerwać bieżącą sesję; wysłać HELLO ponownie; jeśli powtarza się, zgłosić błąd użytkownikowi            |
    | ERR_02 | Format Error | Niepoprawny JSON; brak type lub msg_id; niezgodność typów | Wyświetlić opis błędu; poprawić ramkę i ponowić; w przypadku powtarzających się błędów zakończyć sesję   |
    | ERR_03 | Authorization Timeout   | Brak AUTH w limicie autoryzacji (domyślnie 15 s)  | Wyświetlić timeout; ponowić logowanie; jeśli nadal nieudane, zakończyć sesję   |
    | ERR_04 | Authentication Failed              | Nieprawidłowe dane logowania / API key                                           | Poprosić o poprawne dane; nie próbować automatycznie ponownie bez interwencji użytkownika                |
    | ERR_05 | Token Invalid                      | Token brakujący, wygasły lub niezgodny z Connection ID                           | Wymusić ponowną autoryzację (AUTH); nie akceptować dalszych DATA                                         |
    | ERR_06 | HMAC Mismatch                      | HMAC nie zgadza się z payload                                                    | Odrzucić wiadomość; zalecić ponowne wysłanie z poprawnym HMAC                                            |
    | ERR_07 | Replay Detected                    | Powtórne użycie msg_id/Nonce poza dozwolonym oknem                               | Odrzucić ramkę; zalogować incydent; ewentualnie zablokować sesję przy podejrzeniu ataku                  |
    | ERR_08 | Rate Limit Exceeded                | Przekroczenie limitów (np. >10/s lub >100/min z IP)                              | Tymczasowo zablokować nadawcę; zwrócić info o czasie blokady; host powinien zastosować backoff           |
    | ERR_09 | Delivery Failed Recipient Offline  | Odbiorca niedostępny i brak możliwości dostarczenia                              | Powiadomić nadawcę; nie przechowywać wiadomości; zasugerować ponowne wysłanie po wznowieniu sesji        |
    | ERR_10 | Session Timeout Idle               | Sesja uznana za offline po dłuższym bezruchu                                     | Nadawca powinien spróbować wznowić połączenie; jeśli chce, wykonać pełny handshake                       |
    | ERR_11 | Resource Exhaustion                | Brak zasobów serwera (np. DDoS, brak pamięci)                                    | Zwrócić tymczasowy błąd; nadawca powinien zastosować backoff i retry z jitterem                          |
    | ERR_12 | Internal Server Error              | Błąd wewnętrzny serwera                                                          | Zgłosić użytkownikowi; retry po losowym opóźnieniu                                                       |
    | ERR_13 | Version Mismatch                   | Nieobsługiwana wersja protokołu                                                  | Zaktualizować hosta lub negocjować kompatybilność; zakończyć sesję jeśli niezgodne                       |
    | ERR_14 | Unsupported Media                  | Payload typu/rozmiaru nieobsługiwanego przez protokół                            | Odrzucić; poinformować hosta o dozwolonych typach/limitach                                               |
    | ERR_15 | Not Authorized                     | Brak uprawnień do wykonania akcji                                                | Poprosić o inne uprawnienia; zakończyć próbę wysyłki                                                     |



- Zachowanie po błędach składni/protokołu:
  - W przypadku wykrycia błędu formatu odbiorca natychmiast przerywa przetwarzanie wadliwej ramki i odsyła komunikat `ERROR` z odpowiednim opisem.
  - Poważne naruszenia struktury lub powtarzające się błędy skutkują jednostronnym zamknięciem sesji przez odbiorcę w celu ochrony zasobów.
- Timeouty połączenia:
  - Timeout autoryzacji: Proces `PURR_AUTH` musi zostać sfinalizowany w ciągu 15 sekund od zainicjowania sesji (HELLO).
  - Timeout bezczynności: Jeśli serwer nie odnotuje aktywności (brak `DATA `lub `PING`) przez 60 sekund, klient uznawany jest za offline, a sesja zostaje zamknięta.
- Utrata połączenia w trakcie sesji
  - Dzięki zastosowaniu warstwy transportowej QUIC, KittyProtocol zapewnia wysoką odporność na krótkotrwałe przerwy w transmisji oraz migrację między sieciami (np. Wi-Fi na LTE) bez przerywania sesji.
  - W przypadku trwałej utraty połączenia, protokół przewiduje mechanizm rekonneksji, który umożliwia wznowienie sesji na podstawie ostatniego pomyślnie przetworzonego identyfikatora msg_id.
- Duplikaty i niekompletne wiadomości
  - Duplikaty: Przed atakami typu Replay oraz przypadkowym powieleniem wiadomości (np. przez retransmisję aplikacyjną) chroni pole msg_id, które musi być unikalne w ramach danej sesji.
  - Niekompletność: Integralność strumienia danych gwarantuje warstwa QUIC. Na poziomie aplikacji, każda ramka JSON, która nie przejdzie pełnej walidacji strukturalnej, jest odrzucana jako błąd formatu ERR_02.
- Limity i ochrona przed nadużyciami
  - Rate Limiting: Wprowadzono limit 10 komunikatów na sekundę dla jednego użytkownika, aby zapobiegać spamowaniu i przeciążeniu węzła routingowego.
  - Ochrona IP: Serwer/hosty blokują ruch przekraczający 100 wiadomości na minutę z jednego adresu IP, co stanowi podstawową barierę przed atakami typu DoS.
  - Rozmiar komunikatu: Protokół narzuca maksymalny rozmiar pojedynczej ramki JSON (np. 64 KB), aby zoptymalizować wydajność parsowania i zapobiec atakom polegającym na wysyłaniu nienaturalnie dużych ładunków danych.

- W KittyProtocol błędy są komunikowane poprzez ramki typu `ERROR`, zawierające pola:

  - `type: "ERROR"`
  - `Msg_id`
  - `code` – kod błędu/krótki opis typu błędu
  - `desc` – opis błędu
  
<br><br>
___
# 7. Scenariusze (Przykłady ramek)

   - Scenariusz 1: Poprawna sesja i wymiana wiadomości
   Jest to standardowy przepływ od nawiązania połączenia do przesłania krótkiej wiadomości tekstowej.
       1. Nawiązanie połączenia (Handshake):
       - Nadawca: `{"type": "HELLO", "msg_id": 1, …}`
       - Odbiorca: `{"type": "MEOW_OK", "msg_id": 1, …}`

       2. Autoryzacja (PURR_AUTH):
       - Nadawca: 
      ` {"type": "AUTH", "msg_id": 2, "login": "user_name", "pass": "secret_purr"}`
       - Odbiorca: 
      ` {"type": "MEOW_OK", "msg_id": 2, "token": "abc-123-xyz", "expires": 3600}`

       3. Wymiana danych (DATA):
       - Nadawca: 
       `{"type": "DATA", "msg_id": 3, "token": "abc-123-xyz", "to": "mordunia2", "payload": "Meow! Zobacz tego ASCII kota:", "hmac": "f2d9..."}`
       - Odbiorca (do nadawcy): 
       `{"type": "MEOW_OK", "msg_id": 3, "status": "Delivered"}`

  - Scenariusz 2: Błąd walidacji formatu (ERR_02)
      Scenariusz, w którym nadawca wysyła wiadomość niezgodną ze specyfikacją (np. błędna długość lub brak wymaganych pól).
      - Nadawca: 
      `{"type": "DATA", "msg_id": 10, "token": "abc-123-xyz", "payload": "Zbyt dlugi tekst...", "length": 5}` (Deklarowana długość 5 nie zgadza się z faktyczną)
      - Odbiorca: 
      `{"type": "ERROR", "code": "ERR_02", "msg_id": 10, "desc": "Niepoprawny format danych lub HMAC"}`

   - Scenariusz 3: Timeout procesu autoryzacji (ERR_03)
        Scenariusz obrazujący rygorystyczne limity czasowe protokołu.
        - Nadawca: `{"type": "HELLO", "msg_id": 50}`
        - Odbiorca: `{"type": "MEOW_OK", "msg_id": 50}`
        (Brak aktywności hosta przez ponad 15 sekund)
        - Odbiorca: `{"type": "ERROR", "code": "ERR_03", "desc": "PURR_AUTH timeout reached"}`
        - Połączenie: Odbiorca jednostronnie zamyka sesję QUIC.