package ui_cli

import (
	"bufio"
	"os"

	"github.com/gabbla05/KittyProtocol/client/api"
)

// CliUI implements the UI interface required by the App layer.
// It provides synchronous, blocking terminal input/output.
type CliUI struct {
	client *api.KittyClient
	reader *bufio.Reader
}

// NewCliUI constructs a new CLI frontend bound to a KittyClient instance.
// The UI remains transport-agnostic; the client is used only for ACK callbacks.
func NewCliUI(c *api.KittyClient) *CliUI {
	return &CliUI{
		client: c,
		reader: bufio.NewReader(os.Stdin),
	}
}
