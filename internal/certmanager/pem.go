package certmanager

import (
	"encoding/pem"
	"os"
)

// writePEM writes a PEM-encoded block to the specified file path.
func writePEM(path, pemType string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return pem.Encode(f, &pem.Block{Type: pemType, Bytes: data})
}
