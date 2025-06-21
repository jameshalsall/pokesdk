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

func TestGenerationAPI_List(t *testing.T) {
	ctx := context.Background()

	t.Run("it returns a paginator for listing generations", func(t *testing.T) {
		client, mocks := newGenerationApiForTests(t)
		mocks.backend.HydrateWith([]byte(`{"count": 1, "next": null, "previous": null, "results": [{"name": "generation-i", "url": "https://pokeapi.co/api/v2/generation/1/"}]}`))

		mocks.backend.On("Process", ctx, "http://example.com/generation", map[string]string(nil), mock.Anything).Return(nil).Once()

		pages := client.List()

		firstPage := pages.Next(ctx)
		require.NotNil(t, firstPage)
		require.NoError(t, firstPage.Error)
		require.Len(t, firstPage.Result.Results, 1)
		assert.Equal(t, "generation-i", firstPage.Result.Results[0].Name)
		assert.Equal(t, "https://pokeapi.co/api/v2/generation/1/", firstPage.Result.Results[0].URL)
	})

	t.Run("it should return an error on the page if listing fails", func(t *testing.T) {
		client, mocks := newGenerationApiForTests(t)
		mocks.backend.On("Process", ctx, "http://example.com/generation", map[string]string(nil), mock.Anything).Return(assert.AnError).Once()

		firstPage := client.List().Next(ctx)

		require.NotNil(t, firstPage)
		require.Error(t, firstPage.Error)
		assert.Contains(t, firstPage.Error.Error(), "pokesdk: error listing generations")
	})

	t.Run("it should return an error on the page if the response is malformed", func(t *testing.T) {
		client, mocks := newGenerationApiForTests(t)
		mocks.backend.HydrateWith([]byte(`{"count": 1, "results": [{{{{{{`))

		mocks.backend.On("Process", ctx, "http://example.com/generation", map[string]string(nil), mock.Anything).Return(nil).Once()

		firstPage := client.List().Next(ctx)

		require.NotNil(t, firstPage)
		require.Error(t, firstPage.Error)
		assert.Contains(t, firstPage.Error.Error(), "pokesdk: error listing generations")
	})
}

func TestGenerationAPI_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("it should get a specific generation by ID", func(t *testing.T) {
		client, mocks := newGenerationApiForTests(t)
		mocks.backend.HydrateWith([]byte(`{"id": 1, "name": "generation-i"}`))

		mocks.backend.On("Process", ctx, "http://example.com/generation/1", map[string]string(nil), mock.Anything).Return(nil).Once()

		generation, err := client.GetByID(ctx, 1)

		require.NoError(t, err)
		assert.NotNil(t, generation)
		assert.Equal(t, 1, generation.ID)
		assert.Equal(t, "generation-i", generation.Name)
	})

	t.Run("it should return an error if getting by ID fails", func(t *testing.T) {
		client, mocks := newGenerationApiForTests(t)
		mocks.backend.On("Process", ctx, "http://example.com/generation/999", map[string]string(nil), mock.Anything).Return(assert.AnError).Once()

		_, err := client.GetByID(ctx, 999)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "pokesdk: error getting generation")
	})

	t.Run("it should return an error if the response is malformed", func(t *testing.T) {
		client, mocks := newGenerationApiForTests(t)
		mocks.backend.HydrateWith([]byte(`{"id": 1, "name": "generation-i", {{{{{{`))

		mocks.backend.On("Process", ctx, "http://example.com/generation/1", map[string]string(nil), mock.Anything).Return(nil).Once()

		_, err := client.GetByID(ctx, 1)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "pokesdk: error getting generation")
	})

	t.Run("it should return an error if the generation does not exist", func(t *testing.T) {
		client, mocks := newGenerationApiForTests(t)

		mocks.backend.On("Process", ctx, "http://example.com/generation/999", map[string]string(nil), mock.Anything).Return(backend.ErrResourceNotFound).Once()

		_, err := client.GetByID(ctx, 999)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrGenerationNotFound)
	})
}

func TestGenerationAPI_GetByName(t *testing.T) {
	ctx := context.Background()

	t.Run("it should get a specific generation by name", func(t *testing.T) {
		client, mocks := newGenerationApiForTests(t)
		mocks.backend.HydrateWith([]byte(`{"id": 1, "name": "generation-i"}`))

		mocks.backend.On("Process", ctx, "http://example.com/generation/generation-i", map[string]string(nil), mock.Anything).Return(nil).Once()

		generation, err := client.GetByName(ctx, "generation-i")

		require.NoError(t, err)
		assert.NotNil(t, generation)
		assert.Equal(t, 1, generation.ID)
		assert.Equal(t, "generation-i", generation.Name)
	})

	t.Run("it should return an error if getting by name fails", func(t *testing.T) {
		client, mocks := newGenerationApiForTests(t)
		mocks.backend.On("Process", ctx, "http://example.com/generation/nonexistent", map[string]string(nil), mock.Anything).Return(assert.AnError).Once()

		_, err := client.GetByName(ctx, "nonexistent")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "pokesdk: error getting generation")
	})

	t.Run("it should return an error if the response is malformed", func(t *testing.T) {
		client, mocks := newGenerationApiForTests(t)
		mocks.backend.HydrateWith([]byte(`{"id": 1, "name": "generation-i", {{{{{{`))

		mocks.backend.On("Process", ctx, "http://example.com/generation/generation-i", map[string]string(nil), mock.Anything).Return(nil).Once()

		_, err := client.GetByName(ctx, "generation-i")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "pokesdk: error getting generation")
	})

	t.Run("it should return an error if the generation does not exist", func(t *testing.T) {
		client, mocks := newGenerationApiForTests(t)

		mocks.backend.On("Process", ctx, "http://example.com/generation/nonexistent", map[string]string(nil), mock.Anything).Return(backend.ErrResourceNotFound).Once()

		_, err := client.GetByName(ctx, "nonexistent")

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrGenerationNotFound)
	})
}

func TestGenerationAPI_GetByRef(t *testing.T) {
	ctx := context.Background()

	t.Run("it should get a specific generation by ref", func(t *testing.T) {
		client, mocks := newGenerationApiForTests(t)
		mocks.backend.HydrateWith([]byte(`{"id": 1, "name": "generation-i"}`))

		mocks.backend.On("Process", ctx, "http://example.com/generation/1", map[string]string(nil), mock.Anything).Return(nil).Once()

		ref := GenerationRef{URL: "http://example.com/generation/1"}
		generation, err := client.GetByRef(ctx, ref)

		require.NoError(t, err)
		assert.NotNil(t, generation)
		assert.Equal(t, 1, generation.ID)
		assert.Equal(t, "generation-i", generation.Name)
	})

	t.Run("it should return an error if getting by ref fails", func(t *testing.T) {
		client, mocks := newGenerationApiForTests(t)
		mocks.backend.On("Process", ctx, "http://example.com/generation/1", map[string]string(nil), mock.Anything).Return(assert.AnError).Once()

		ref := GenerationRef{URL: "http://example.com/generation/1"}
		_, err := client.GetByRef(ctx, ref)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "pokesdk: error getting generation")
	})

	t.Run("it should return an error if the response is malformed", func(t *testing.T) {
		client, mocks := newGenerationApiForTests(t)
		mocks.backend.HydrateWith([]byte(`{"id": 1, "name": "generation-i", {{{{{{`))

		mocks.backend.On("Process", ctx, "http://example.com/generation/generation-i", map[string]string(nil), mock.Anything).Return(nil).Once()

		ref := GenerationRef{URL: "http://example.com/generation/1"}
		_, err := client.GetByRef(ctx, ref)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "pokesdk: error getting generation")
	})

	t.Run("it should return an error if the generation does not exist", func(t *testing.T) {
		client, mocks := newGenerationApiForTests(t)

		mocks.backend.On("Process", ctx, "http://example.com/generation/99", map[string]string(nil), mock.Anything).Return(backend.ErrResourceNotFound).Once()

		ref := GenerationRef{URL: "http://example.com/generation/99"}
		_, err := client.GetByRef(ctx, ref)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrGenerationNotFound)
	})
}

type generationApiMocks struct {
	backend *pokesdktest.MockBackend
}

func newGenerationApiForTests(t *testing.T) (GenerationAPI, generationApiMocks) {
	t.Helper()

	mocks := generationApiMocks{
		backend: &pokesdktest.MockBackend{},
	}

	return GenerationAPI{cfg: Config{baseURL: "http://example.com/", backend: mocks.backend}}, mocks
}
