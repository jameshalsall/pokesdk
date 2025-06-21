package pokesdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/jameshalsall/pokesdk/internal/backend"
	"github.com/jameshalsall/pokesdk/pokesdktest"
)

func TestPokemonAPI_List(t *testing.T) {
	ctx := context.Background()

	t.Run("it returns a paginator for listing pokemon", func(t *testing.T) {
		client, mocks := newPokemonApiForTests(t)
		mocks.backend.HydrateWith([]byte(`{"count": 1, "next": null, "previous": null, "results": [{"name": "bulbasaur", "url": "https://pokeapi.co/api/v2/pokemon/1/"}]}`))

		mocks.backend.On("Process", ctx, "http://example.com/pokemon", map[string]string(nil), mock.Anything).Return(nil).Once()

		pages := client.List()

		firstPage := pages.Next(ctx)
		require.NotNil(t, firstPage)
		require.NoError(t, firstPage.Error)
		require.Len(t, firstPage.Result.Results, 1)
		assert.Equal(t, "bulbasaur", firstPage.Result.Results[0].Name)
		assert.Equal(t, "https://pokeapi.co/api/v2/pokemon/1/", firstPage.Result.Results[0].URL)
	})

	t.Run("it should return an error on the page if listing fails", func(t *testing.T) {
		client, mocks := newPokemonApiForTests(t)
		mocks.backend.On("Process", ctx, "http://example.com/pokemon", map[string]string(nil), mock.Anything).Return(assert.AnError).Once()

		firstPage := client.List().Next(ctx)

		require.NotNil(t, firstPage)
		require.Error(t, firstPage.Error)
		assert.Contains(t, firstPage.Error.Error(), "pokesdk: error listing pokemon")
	})

	t.Run("it should return an error on the page if the response is malformed", func(t *testing.T) {
		client, mocks := newPokemonApiForTests(t)
		mocks.backend.HydrateWith([]byte(`{"count": 1, "results": [{{{{{{`))

		mocks.backend.On("Process", ctx, "http://example.com/pokemon", map[string]string(nil), mock.Anything).Return(nil).Once()

		firstPage := client.List().Next(ctx)

		require.NotNil(t, firstPage)
		require.Error(t, firstPage.Error)
		assert.Contains(t, firstPage.Error.Error(), "pokesdk: error listing pokemon")
	})
}

func TestPokemonAPI_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("it should get a specific pokemon by ID", func(t *testing.T) {
		client, mocks := newPokemonApiForTests(t)
		mocks.backend.HydrateWith([]byte(`{"id": 1, "name": "bulbasaur"}`))

		mocks.backend.On("Process", ctx, "http://example.com/pokemon/1", map[string]string(nil), mock.Anything).Return(nil).Once()

		pokemon, err := client.GetByID(ctx, 1)

		require.NoError(t, err)
		assert.NotNil(t, pokemon)
		assert.Equal(t, 1, pokemon.ID)
		assert.Equal(t, "bulbasaur", pokemon.Name)
	})

	t.Run("it should return an error if getting by ID fails", func(t *testing.T) {
		client, mocks := newPokemonApiForTests(t)
		mocks.backend.On("Process", ctx, "http://example.com/pokemon/999", map[string]string(nil), mock.Anything).Return(assert.AnError).Once()

		_, err := client.GetByID(ctx, 999)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "pokesdk: error getting pokemon")
	})

	t.Run("it should return an error if the response is malformed", func(t *testing.T) {
		client, mocks := newPokemonApiForTests(t)
		mocks.backend.HydrateWith([]byte(`{"id": 1, "name": "bulbasaur", "abilities": [{{{{{{`))

		mocks.backend.On("Process", ctx, "http://example.com/pokemon/1", map[string]string(nil), mock.Anything).Return(nil).Once()

		_, err := client.GetByID(ctx, 1)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "pokesdk: error getting pokemon")
	})

	t.Run("it should return an error if the pokemon does not exist", func(t *testing.T) {
		client, mocks := newPokemonApiForTests(t)

		mocks.backend.On("Process", ctx, "http://example.com/pokemon/999", map[string]string(nil), mock.Anything).Return(backend.ErrResourceNotFound).Once()

		_, err := client.GetByID(ctx, 999)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrPokemonNotFound)
	})
}

func TestPokemonAPI_GetByName(t *testing.T) {
	ctx := context.Background()

	t.Run("it should get a specific pokemon by name", func(t *testing.T) {
		client, mocks := newPokemonApiForTests(t)
		mocks.backend.HydrateWith([]byte(`{"id": 1, "name": "bulbasaur"}`))

		mocks.backend.On("Process", ctx, "http://example.com/pokemon/bulbasaur", map[string]string(nil), mock.Anything).Return(nil).Once()

		pokemon, err := client.GetByName(ctx, "bulbasaur")

		require.NoError(t, err)
		assert.NotNil(t, pokemon)
		assert.Equal(t, 1, pokemon.ID)
		assert.Equal(t, "bulbasaur", pokemon.Name)
	})

	t.Run("it should return an error if getting by name fails", func(t *testing.T) {
		client, mocks := newPokemonApiForTests(t)
		mocks.backend.On("Process", ctx, "http://example.com/pokemon/nonexistent", map[string]string(nil), mock.Anything).Return(assert.AnError).Once()

		_, err := client.GetByName(ctx, "nonexistent")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "pokesdk: error getting pokemon")
	})

	t.Run("it should return an error if the response is malformed", func(t *testing.T) {
		client, mocks := newPokemonApiForTests(t)
		mocks.backend.HydrateWith([]byte(`{"id": 1, "name": "bulbasaur", "abilities": [{{{{{{`))

		mocks.backend.On("Process", ctx, "http://example.com/pokemon/bulbasaur", map[string]string(nil), mock.Anything).Return(nil).Once()

		_, err := client.GetByName(ctx, "bulbasaur")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "pokesdk: error getting pokemon")
	})

	t.Run("it should return an error if the pokemon does not exist", func(t *testing.T) {
		client, mocks := newPokemonApiForTests(t)

		mocks.backend.On("Process", ctx, "http://example.com/pokemon/nonexistent", map[string]string(nil), mock.Anything).Return(backend.ErrResourceNotFound).Once()

		_, err := client.GetByName(ctx, "nonexistent")

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrPokemonNotFound)
	})
}

func TestPokemonAPI_GetByRef(t *testing.T) {
	ctx := context.Background()

	t.Run("it should get a specific pokemon by ref", func(t *testing.T) {
		client, mocks := newPokemonApiForTests(t)
		mocks.backend.HydrateWith([]byte(`{"id": 1, "name": "bulbasaur"}`))

		mocks.backend.On("Process", ctx, "http://example.com/pokemon/1", map[string]string(nil), mock.Anything).Return(nil).Once()

		ref := PokemonRef{URL: "http://example.com/pokemon/1"}
		pokemon, err := client.GetByRef(ctx, ref)

		require.NoError(t, err)
		assert.NotNil(t, pokemon)
		assert.Equal(t, 1, pokemon.ID)
		assert.Equal(t, "bulbasaur", pokemon.Name)
	})

	t.Run("it should return an error if getting by ref fails", func(t *testing.T) {
		client, mocks := newPokemonApiForTests(t)
		mocks.backend.On("Process", ctx, "http://example.com/pokemon/99", map[string]string(nil), mock.Anything).Return(assert.AnError).Once()

		ref := PokemonRef{URL: "http://example.com/pokemon/99"}
		_, err := client.GetByRef(ctx, ref)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "pokesdk: error getting pokemon")
	})

	t.Run("it should return an error if the response is malformed", func(t *testing.T) {
		client, mocks := newPokemonApiForTests(t)
		mocks.backend.HydrateWith([]byte(`{"id": 1, "name": "bulbasaur", "abilities": [{{{{{{`))

		mocks.backend.On("Process", ctx, "http://example.com/pokemon/bulbasaur", map[string]string(nil), mock.Anything).Return(nil).Once()

		_, err := client.GetByName(ctx, "bulbasaur")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "pokesdk: error getting pokemon")
	})

	t.Run("it should return an error if the pokemon does not exist", func(t *testing.T) {
		client, mocks := newPokemonApiForTests(t)

		mocks.backend.On("Process", ctx, "http://example.com/pokemon/99", map[string]string(nil), mock.Anything).Return(backend.ErrResourceNotFound).Once()

		ref := PokemonRef{URL: "http://example.com/pokemon/99"}
		_, err := client.GetByRef(ctx, ref)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrPokemonNotFound)
	})
}

type pokemonApiMocks struct {
	backend *pokesdktest.MockBackend
}

func newPokemonApiForTests(t *testing.T) (PokemonAPI, pokemonApiMocks) {
	t.Helper()

	mocks := pokemonApiMocks{
		backend: &pokesdktest.MockBackend{},
	}

	return PokemonAPI{cfg: Config{baseURL: "http://example.com/", backend: mocks.backend}}, mocks
}
