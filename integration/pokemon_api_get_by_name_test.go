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

func TestPokemonGetByName(t *testing.T) {
	server := NewMockServer(t)
	defer server.Close()

	ctx := context.Background()

	server.StubGET("/pokemon/venusaur", Response{
		StatusCode: 200,
		Body:       pokemonResponse,
	})

	t.Run("it returns a specific pokemon by name", func(t *testing.T) {
		client := server.PokeSDKClient()

		pokemon, err := client.Pokemon.GetByName(ctx, "venusaur")
		require.NoError(t, err)

		assert.NotNil(t, pokemon)
		assert.Equal(t, "venusaur", pokemon.Name)
		assert.Equal(t, 3, pokemon.ID)

		reqs := server.Requests()
		require.Len(t, reqs, 1)
		assert.Equal(t, "/pokemon/venusaur", reqs[0].Path)
	})

	t.Run("it returns an error if pokemon not found", func(t *testing.T) {
		server.Reset()
		client := server.PokeSDKClient()

		_, err := client.Pokemon.GetByName(ctx, "nonexistent")
		require.Error(t, err)
		assert.ErrorIs(t, err, pokesdk.ErrPokemonNotFound)

		reqs := server.Requests()
		require.Len(t, reqs, 1)
		assert.Equal(t, "/pokemon/nonexistent", reqs[0].Path)
	})
}
