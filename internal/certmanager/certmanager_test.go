package certmanager

import (
	"crypto/tls"
	"os"
	"testing"
)

func TestSetupTLSConfig(t *testing.T) {
	tempCert := "temp_cert.pem"
	tempKey := "temp_key.pem"

	// Czyszczenie po teście
	defer os.Remove(tempCert)
	defer os.Remove(tempKey)

	// Etap 1: Wygenerowanie nowych certyfikatów
	config, err := SetupTLSConfig(tempCert, tempKey)
	if err != nil {
		t.Fatalf("Oczekiwano sukcesu, otrzymano błąd: %v", err)
	}

	// Sprawdzenie, czy wymuszono TLS 1.3
	if config.MinVersion != tls.VersionTLS13 {
		t.Errorf("Oczekiwano TLS 1.3, otrzymano: %x", config.MinVersion)
	}

	// Sprawdzenie, czy pliki faktycznie powstały
	if _, err := os.Stat(tempCert); os.IsNotExist(err) {
		t.Errorf("Plik certyfikatu nie został wygenerowany")
	}

	// Etap 2: Wczytanie z już istniejących plików (bez generowania)
	config2, err := SetupTLSConfig(tempCert, tempKey)
	if err != nil {
		t.Fatalf("Oczekiwano sukcesu przy wczytywaniu istniejących plików, błąd: %v", err)
	}

	if len(config2.Certificates) == 0 {
		t.Errorf("Brak certyfikatów w konfiguracji")
	}
}
