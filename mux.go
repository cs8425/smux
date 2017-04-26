package smux

import (
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
)

// Config is used to tune the Smux session
type Config struct {
	// KeepAliveInterval is how often to send a NOP command to the remote
	KeepAliveInterval time.Duration

	// KeepAliveTimeout is how long the session
	// will be closed if no data has arrived
	KeepAliveTimeout time.Duration

	// MaxFrameSize is used to control the maximum
	// frame size to sent to the remote
	MaxFrameSize int

	// MaxReceiveBuffer is used to control the maximum
	// number of data in the buffer pool
	MaxReceiveBuffer int

	// Enable Stream buffer
	EnableStreamBuffer bool

	// maximum bytes that each Stream can use
	MaxStreamBuffer int

	// send cmdEMP earlier when remain ?? bytes data in buffer
	MinStreamBuffer int
}

// DefaultConfig is used to return a default configuration
func DefaultConfig() *Config {
	return &Config{
		KeepAliveInterval:  10 * time.Second,
		KeepAliveTimeout:   30 * time.Second,
		MaxFrameSize:       4096,
		MaxReceiveBuffer:   16 * 1024 * 1024,
		EnableStreamBuffer: false,
		MaxStreamBuffer:    16384,
		MinStreamBuffer:    4096,
	}
}

// VerifyConfig is used to verify the sanity of configuration
func VerifyConfig(config *Config) error {
	if config.KeepAliveInterval == 0 {
		return errors.New("keep-alive interval must be positive")
	}
	if config.KeepAliveTimeout < config.KeepAliveInterval {
		return fmt.Errorf("keep-alive timeout must be larger than keep-alive interval")
	}
	if config.MaxFrameSize <= 0 {
		return errors.New("max frame size must be positive")
	}
	if config.MaxFrameSize > 65535 {
		return errors.New("max frame size must not be larger than 65535")
	}
	if config.MaxReceiveBuffer <= 0 {
		return errors.New("max receive buffer must be positive")
	}
	if config.MaxStreamBuffer <= 0 {
		return errors.New("max stream receive buffer must be positive")
	}
	if config.MinStreamBuffer < 0 || config.MinStreamBuffer > config.MaxStreamBuffer {
		return errors.New("min stream receive buffer must be positive and <= MaxStreamBuffer")
	}
	return nil
}

// Server is used to initialize a new server-side connection.
func Server(conn io.ReadWriteCloser, config *Config) (*Session, error) {
	if config == nil {
		config = DefaultConfig()
	}
	if err := VerifyConfig(config); err != nil {
		return nil, err
	}
	return newSession(config, conn, false), nil
}

// Client is used to initialize a new client-side connection.
func Client(conn io.ReadWriteCloser, config *Config) (*Session, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if err := VerifyConfig(config); err != nil {
		return nil, err
	}
	return newSession(config, conn, true), nil
}
