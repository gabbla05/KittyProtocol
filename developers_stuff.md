hub:  go run hub/main.go
clients:  go run client/main.go

<img width="1342" height="493" alt="image" src="https://github.com/user-attachments/assets/0d654560-ff16-4aa3-92a5-66d7defae7ba" />


TODO:
Potwierdzenia dostarczenia (Delivery ACK): Wysyłanie ramki MEOW_OK przez odbiorcę natychmiast po poprawnym odebraniu i sparsowaniu ramki DATA. Pozwoli to nadawcy na zmianę statusu wiadomości w konsoli na "Delivered".

Obsługa błędów formatu (ERR_02): Rygorystyczna walidacja przychodzących danych – jeśli JSON jest niekompletny lub ma błędne typy pól, aplikacja musi odesłać ERROR z kodem ERR_02 i przerwać przetwarzanie.

Logowanie zdarzeń (Diagnostic Logs): Mchanizm zapisywania każdej wysłanej i odebranej ramki do pliku tekstowego lub konsoli z naniesionym znacznikiem czasu. Ułatwi to debugowanie i będzie wymagane jako element diagnostyki w Etapie 2.

Implementacja Timeoutów: 15-sekundowy limit na przesłanie ramki AUTH po nawiązaniu połączenia HELLO. W przypadku braku reakcji, serwer musi odesłać ERR_03 i zamknąć sesję.

Mechanizm Keep-alive (PING): Automatyczne wysyłanie rami PING co 30 sekund w trakcie bezczynności, aby zapobiec zamykaniu strumienia przez urządzenia sieciowe. Odbiorca powinien odpowiadać na PING zwykłym MEOW_OK.

Dokumentacja Etapu 2: Spisywanie przypadków użycia (use cases), opisując m.in. udane logowanie, wysyłkę P2P oraz reakcję na offline rozmówcy. Muszą udowadniać, że aplikacja realizuje założenia protokołu.