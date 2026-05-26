package secretstore

import (
	"encoding/base64"
	"encoding/json"
)

// diskSecrets is the on‑disk representation of secrets.
// Each secret is base64‑encoded to keep the JSON printable.
type diskSecrets struct {
	Peers map[string]string `json:"peers"` // peer -> base64(secret)
}

// marshalDisk converts the in‑memory map into JSON plaintext.
func marshalDisk(secrets map[string][]byte) ([]byte, error) {
	ds := diskSecrets{
		Peers: make(map[string]string, len(secrets)),
	}
	for peer, secret := range secrets {
		ds.Peers[peer] = base64.StdEncoding.EncodeToString(secret)
	}
	return json.MarshalIndent(ds, "", "  ")
}

// unmarshalDisk parses JSON plaintext into an in‑memory map.
func unmarshalDisk(data []byte) (map[string][]byte, error) {
	var ds diskSecrets
	if err := json.Unmarshal(data, &ds); err != nil {
		return nil, err
	}

	out := make(map[string][]byte, len(ds.Peers))
	for peer, b64 := range ds.Peers {
		raw, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			continue // skip malformed entries
		}
		out[peer] = raw
	}
	return out, nil
}
