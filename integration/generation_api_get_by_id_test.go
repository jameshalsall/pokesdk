//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jameshalsall/pokesdk"
)

func TestGenerationGetByID(t *testing.T) {
	server := NewMockServer(t)
	defer server.Close()

	ctx := context.Background()

	t.Run("it returns a specific generation by ID", func(t *testing.T) {
		server.StubGET("/generation/1", Response{
			StatusCode: 200,
			Body:       generationResponse,
		})

		client := server.PokeSDKClient()

		generation, err := client.Generation.GetByID(ctx, 1)
		require.NoError(t, err)

		assert.NotNil(t, generation)
		assert.Equal(t, "generation-i", generation.Name)
		assert.Equal(t, 1, generation.ID)

		reqs := server.Requests()
		require.Len(t, reqs, 1)
		assert.Equal(t, "/generation/1", reqs[0].Path)
	})

	t.Run("it returns an error if generation not found", func(t *testing.T) {
		server.Reset()
		client := server.PokeSDKClient()

		_, err := client.Generation.GetByID(ctx, 99)
		require.Error(t, err)
		assert.ErrorIs(t, err, pokesdk.ErrGenerationNotFound)

		reqs := server.Requests()
		require.Len(t, reqs, 1)
		assert.Equal(t, "/generation/99", reqs[0].Path)
	})
}
