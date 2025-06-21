//go:build integration

package integration

import (
	"context"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jameshalsall/pokesdk"
)

func TestGenerationGetByName(t *testing.T) {
	server := NewMockServer(t)
	defer server.Close()

	ctx := context.Background()

	server.StubGET("/generation/generation-i", Response{
		StatusCode: 200,
		Body:       generationResponse,
	})

	t.Run("it returns a specific generation by name", func(t *testing.T) {
		client := server.PokeSDKClient()

		generation, err := client.Generation.GetByName(ctx, "generation-i")
		require.NoError(t, err)

		assert.NotNil(t, generation)
		assert.Equal(t, "generation-i", generation.Name)
		assert.Equal(t, 1, generation.ID)

		reqs := server.Requests()
		require.Len(t, reqs, 1)
		assert.Equal(t, "/generation/generation-i", reqs[0].Path)
	})

	t.Run("it returns an error if generation not found", func(t *testing.T) {
		server.Reset()
		client := server.PokeSDKClient()

		_, err := client.Generation.GetByName(ctx, "nonexistent")
		require.Error(t, err)
		assert.ErrorIs(t, err, pokesdk.ErrGenerationNotFound)

		reqs := server.Requests()
		require.Len(t, reqs, 1)
		assert.Equal(t, "/generation/nonexistent", reqs[0].Path)
	})
}
