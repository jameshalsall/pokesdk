package pokesdk

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jameshalsall/pokesdk/internal/backend"
)

func TestNewConfig(t *testing.T) {
	t.Run("it has sensible defaults", func(t *testing.T) {
		cfg := NewConfig()

		assert.NotNil(t, cfg.backend)

		if _, ok := cfg.backend.(*backend.HTTP); !ok {
			t.Fatalf("expected backend to be of type *HTTP, got %T", cfg.backend)
		}

		assert.Equal(t, defaultBaseAPIURL, cfg.baseURL)
	})

	t.Run("it can be configured with a custom base URL", func(t *testing.T) {
		customBaseURL := "https://custom.pokeapi.co/api/v2"
		cfg := NewConfig(WithCustomBaseURL(customBaseURL))

		assert.Equal(t, customBaseURL, cfg.baseURL)
	})

	t.Run("it can be configured with a custom HTTP client", func(t *testing.T) {
		client := &http.Client{}

		cfg := NewConfig(WithCustomHttpClient(client))

		_, ok := cfg.backend.(*backend.HTTP)
		assert.True(t, ok)
	})
}
