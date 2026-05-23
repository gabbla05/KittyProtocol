  🇵🇱 *Wersja polska* | [🇬🇧 English version](README.md)
 
 # Kitty Protocol

![image](resources/img/kitty_logo.png)

 ## Opis projektu

 Kitty Protocol to autorski, lekki protokół komunikacyjny zaprojektowany w celu zapewnienia maksymalnej prywatności i szybkości wymiany krótkich wiadomości tekstowych. Projekt opiera się na architekturze Klient-Serwer, gdzie centralny serwer (Hub) pełni jedynie rolę routera, nie mając dostępu do treści rozmów.

 ### Kluczowe cechy:
 * **`Bezpieczeństwo End-to-End (E2EE)`**: Pełne szyfrowanie realizowane wyłącznie po stronie klientów; serwer nie przechowuje kluczy ani historii wiadomości.
 * **`Wykorzystanie QUIC`**: Protokół bazuje na warstwie transportowej QUIC, co gwarantuje niskie opóźnienia (0-RTT) oraz stabilność połączenia przy zmianie sieci (np. z Wi-Fi na LTE).
 * **`Format JSON`**: Dane przesyłane są w czytelnym formacie tekstowym JSON, co ułatwia debugowanie i rozwój protokołu.
 * **`Odporność na ataki`**: Natywne wsparcie dla TLS 1.3 oraz mechanizmy zapobiegające atakom typu Man-in-the-Middle (MitM) i Replay.

 ### Dokumentacja:
 Szczegółowa dokumentacja protokołu znajduje się w dokumencie: [KittyProtocol.pdf](documentation/KittyProtocol.pdf)

 ---
 Authors: 
 
 [Gabriela Błaut](https://github.com/gabbla05), 

 [Michał Brzeziński](https://github.com/Michal-Brzezinski), 
 
 [Aleksandra Gołek](https://github.com/styliana)  
