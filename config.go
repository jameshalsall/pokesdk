package pokesdk

import (
	"context"

	"github.com/jameshalsall/pokesdk/internal/backend"
)

const (
	defaultBaseAPIURL = "https://pokeapi.co/api/v2"
)

// Option defines a function that can be used to configure the PokeSDK.
type Option func(cfg *Config)

// Backend defines the interface for a backend that can process requests to the PokeAPI.
// It abstracts the underlying implementation, allowing for different backends to be used in the future.
type Backend interface {
	Process(ctx context.Context, url string, params map[string]string, out any) error
}

type Config struct {
	baseURL string
	backend Backend
}

// NewConfig creates a new Config instance with the provided options applied.
func NewConfig(opts ...Option) Config {
	cfg := Config{}
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.backend == nil {
		WithHttpBackend()(&cfg)
	}

	if cfg.baseURL == "" {
		WithDefaultBaseURL()(&cfg)
	}

	return cfg
}

// WithHttpBackend sets the backend to a default HTTP backend.
func WithHttpBackend() Option {
	return func(cfg *Config) {
		cfg.backend = backend.NewDefaultHTTP()
	}
}

// WithCustomHttpClient sets an HTTP backend in the Config.
// The provided client will be used to create a new HTTP backend.
func WithCustomHttpClient(client backend.HTTPClient) Option {
	return func(cfg *Config) {
		cfg.backend = backend.NewHTTP(client)
	}
}

// WithCustomBaseURL sets a custom base URL in the Config.
func WithCustomBaseURL(baseURL string) Option {
	return func(cfg *Config) {
		cfg.baseURL = baseURL
	}
}

// WithDefaultBaseURL sets the base URL to the default PokeAPI URL.
func WithDefaultBaseURL() Option {
	return func(cfg *Config) {
		cfg.baseURL = defaultBaseAPIURL
	}
}
