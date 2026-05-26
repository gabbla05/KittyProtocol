# Proponowana ścieżka konfiguracji ogólnodostępnego huba korzystając z platformy Azure

1. Wejdź na stronę Azure i załóż konto.
   
    Przy pierwszej rejestracji dostajesz 200$ kredytu na 30 dni – wszystkie poniższe kroki mieszczą się spokojnie w tym limicie.

2. Utworzenie i konfiguracja maszyny wirtualnej
   1. Utwórz maszynę wirtualną z dowolnymi parametrami
   2. Po utworzeniu maszyny wykonaj następujące kroki w celu otwarcia portu do nasłuchu:
      1. W widoku swojej maszyny wirtualnej wejdź w: Sieć -> Ustawienia Sieci
      2. W sekcji Reguły wybierz "Dodaj nową regułę" i wypełnij:
        - Źródło: any
        - Źródłowe zakresy portów: *
        - Lokalizacja docelowa: Any
        - Usługa: Custom
        - Docelowe zakresy portów: \<twój port\>
        - Protokół: udp (protokół quic oparty jest na udp)
        - Priorytet: np. 100
        - nazwa: dowolność
      3. Zapisz konfigurację dodanej reguły
   3. Uruchom ponownie maszynę wirtualną

3. Połączenie z VM przez SSH: 

    ```bash
        ssh -i /ścieżka/do/klucza.pem azureusername@PUBLIC_IP
    ```
    Ścieżka do key.pem najprawdopodoniej będzie czymś w rodzaju:
    
    ```bash
        \<twoje foldery\>.../KittyProtocol/certs/key.pem
    ```
    `PUBLIC_IP` znajdziesz w szczegółach VM.

4. Instalacja Go, Git, Docker na VM
    ```bash
        sudo apt update
        sudo apt install -y golang-go git
        go version
        git --version

        sudo apt update
        sudo apt install -y ca-certificates curl

        sudo install -m 0755 -d /etc/apt/keyrings
        sudo curl -fsSL https://docker.com -o /etc/apt/keyrings/docker.asc
        sudo chmod a+r /etc/apt/keyrings/docker.asc

        echo \
            "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://docker.com \
            $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
            sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

        sudo apt update
        sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    ```

    Dla wygody sugerowane są także działania:

    ```bash
        # 1. Utwórz grupę docker (na wypadek, gdyby instalator jej nie stworzył)
        sudo groupadd docker

        # 2. Dodaj swojego obecnego użytkownika ($USER) do grupy docker
        sudo usermod -aG docker $USER

        # 3. Aktywuj zmiany w obecnej sesji terminala (bez restartu systemu)
        newgrp docker
    ```

    Zlikwiduje to powinność każdorazowego używania `sudo`

    
5. Pobranie repozytorium

    ```bash
        cd ~
        git clone https://github.com/gabbla05/KittyProtocol.git
        cd KittyProtocol
    ```

6. Konfiguracja plików .env:

    wykonaj operacje:

    ```bash
        cd KittyProtocol
        cp ./.env.example ./.env
        cp ./db/.env.example ./db/.env
    ```

7. Uruchomienie kontenerów z usługami bazy danych PostgreSQL

    ```bash
        cd KittyProtocol/db
        docker compose up -d
    ```

8. Wejście do środka kontenera bazy danych:

    Operacja dla przykładowej konfiguracji

    ```bash
        docker compose exec db psql -U kitty -d kittyhub
    ```

9. Utworzenie tabeli użytkowników

    Wykonaj skrypt sql weewnątrz kontenera:

    ```sql
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username TEXT UNIQUE NOT NULL,
            password_hash TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT NOW()
        );
    ```

    Sprawdź czy tabela została utworzona, a następnie wyjdź z kontenera:
    
    ```sql
        \dt --Wyświetlenie listy wszystkich tabel w bazie.
        \q  --Wyjście z konsoli PostgreSQL z powrotem do Twojego systemu.
    ```

10. Ustawienie zmiennych środowiskowych dla Huba

    Na VM, np. w pliku `~/.bashrc` dodaj na końcu:

    ```bash
        export KITTY_DB_DSN="postgres://kitty:kittypass@localhost:5432/kittyhub?sslmode=disable"
        export KITTY_INTERCEPT_ADDR="0.0.0.0:9999"
        export KITTY_LOG_COLOR=0
        export HUB_TLS_CERT="/home/kittyhub/KittyProtocol/certs/cert.pem"
        export HUB_TLS_KEY="/home/kittyhub/KittyProtocol/certs/key.pem"
        export HUB_LOG_LEVEL="info"
    ```

    oraz załaduj:

    ```bash
        source ~/.bashrc
    ```

11. Certyfikaty TLS

    Przy uruchomieniu huba certyfikaty powinny wygenerować się automatycznie, jeżeli na ten moment ich nie ma

12. Uruchomienie Huba

    ```bash    
        cd ~/KittyProtocol
        go run ./cmd/hub
    ```

# Konfiguracja klienta kompatybilnego z powyższą konfiguracją huba

1. Na lokalnym urządzeniu w głównym katalogu projektu należy zmodyfikować plik `.env` zgodnie ze wskazówkami z pliku .env.example:

   ```env

        # -----------------------------
        # KittyProtocol configuration
        # -----------------------------

        # 127.0.0.1:9999 - local version

        # ==========   AZURE   ============
        # for KITTY_HUB_ADDR you can assign
        # 6.67.67.67:9999 - Azure version
        # it depends on which port you
        # enabled in network settings 
        # as a role and assigned public ip

        KITTY_HUB_ADDR=127.0.0.1:9999


        KITTY_INTERCEPT_ADDR=127.0.0.1:9999
        # 127.0.0.1:9999 - version for local dev tests
        # 0.0.0.0:9999 - version for production - listen all devices

   ``` 

## Proces można oczywiście zautomatyzować przy pomocy narzędzi takich jak Terraform i Ansible  