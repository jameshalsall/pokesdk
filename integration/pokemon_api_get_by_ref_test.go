//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jameshalsall/pokesdk"
)

func TestPokemonGetByRef(t *testing.T) {
	server := NewMockServer(t)
	defer server.Close()

	ctx := context.Background()

	t.Run("it returns a specific pokemon using a ref", func(t *testing.T) {
		server.StubGET("/pokemon/3", Response{
			StatusCode: 200,
			Body:       pokemonResponse,
		})

		client := server.PokeSDKClient()

		ref := pokesdk.PokemonRef{
			Name: "venusaur",
			URL:  server.URL() + "/pokemon/3",
		}
		pokemon, err := client.Pokemon.GetByRef(ctx, ref)
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

		ref := pokesdk.PokemonRef{
			Name: "notfound",
			URL:  server.URL() + "/pokemon/99",
		}
		_, err := client.Pokemon.GetByRef(ctx, ref)
		require.Error(t, err)
		assert.ErrorIs(t, err, pokesdk.ErrPokemonNotFound)

		reqs := server.Requests()
		require.Len(t, reqs, 1)
		assert.Equal(t, "/pokemon/99", reqs[0].Path)
	})
}
