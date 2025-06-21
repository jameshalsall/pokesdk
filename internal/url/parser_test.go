package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathFromURL(t *testing.T) {
	t.Run("returns path and query from URL", func(t *testing.T) {
		rawUrl := "https://example.com/api/v1/resource?query=param"

		got, ok := PathFromURL(rawUrl)
		assert.True(t, ok)
		assert.Equal(t, "/api/v1/resource?query=param", got)
	})

	t.Run("returns empty string for invalid URL", func(t *testing.T) {
		rawUrl := ":adwad"

		got, ok := PathFromURL(rawUrl)
		assert.False(t, ok)
		assert.Equal(t, "", got)
	})
}
