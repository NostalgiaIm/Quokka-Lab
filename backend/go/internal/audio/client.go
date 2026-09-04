package audio

import "log/slog"

// Client describes the boundary between Go and the C++ audio engine.
type Client struct {
	socketPath string
	binaryPath string
	logger     *slog.Logger
}

// NewClient prepares the UDS client boundary. The command-line binary remains
// available as a fallback while the Cap'n Proto protocol is being implemented.
func NewClient(socketPath string, binaryPath string, logger *slog.Logger) *Client {
	return &Client{socketPath: socketPath, binaryPath: binaryPath, logger: logger}
}

func (c *Client) SocketPath() string {
	return c.socketPath
}

func (c *Client) BinaryPath() string {
	return c.binaryPath
}
