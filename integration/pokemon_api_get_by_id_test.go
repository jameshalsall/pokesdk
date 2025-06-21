//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jameshalsall/pokesdk"
)

func TestPokemonGetByID(t *testing.T) {
	server := NewMockServer(t)
	defer server.Close()

	ctx := context.Background()

	t.Run("it returns a specific pokemon by ID", func(t *testing.T) {
		server.StubGET("/pokemon/3", Response{
			StatusCode: 200,
			Body:       pokemonResponse,
		})

		client := server.PokeSDKClient()

		pokemon, err := client.Pokemon.GetByID(ctx, 3)
		require.NoError(t, err)

		assert.NotNil(t, pokemon)
		assert.Equal(t, "venusaur", pokemon.Name)
		assert.Equal(t, 3, pokemon.ID)

		reqs := server.Requests()
		require.Len(t, reqs, 1)
		assert.Equal(t, "/pokemon/3", reqs[0].Path)
	})

	t.Run("it returns an error if pokemon not found", func(t *testing.T) {
		server.Reset()
		client := server.PokeSDKClient()

		_, err := client.Pokemon.GetByID(ctx, 99)
		require.Error(t, err)
		assert.ErrorIs(t, err, pokesdk.ErrPokemonNotFound)

		reqs := server.Requests()
		require.Len(t, reqs, 1)
		assert.Equal(t, "/pokemon/99", reqs[0].Path)
	})
}
