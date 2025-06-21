//go:build integration

package integration

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerationAPIList(t *testing.T) {
	ctx := context.Background()
	server := NewMockServer(t)
	defer server.Close()

	server.StubGET("/generation", Response{
		StatusCode: http.StatusOK,
		Body:       ResponseBytes(listGenerationResponsePage1, server.URL()),
		Headers:    map[string]string{"Content-Type": "application/json"},
	})

	server.StubGET("/generation?offset=20&limit=20", Response{
		StatusCode: http.StatusOK,
		Body:       ResponseBytes(listGenerationResponsePage2, server.URL()),
		Headers:    map[string]string{"Content-Type": "application/json"},
	})

	t.Run("it lists pages of generations", func(t *testing.T) {
		client := server.PokeSDKClient()
		pages := client.Generation.List()

		firstPage := pages.Next(ctx)
		require.NotNil(t, firstPage)
		require.Len(t, firstPage.Result.Results, 1)
		assert.Equal(t, "generation-i", firstPage.Result.Results[0].Name)
		assert.Equal(t, "https://pokeapi.co/api/v2/generation/1/", firstPage.Result.Results[0].URL)
		assert.Equal(t, 1, firstPage.Result.Count)

		secondPage := pages.Next(ctx)
		require.NotNil(t, secondPage)
		require.Len(t, secondPage.Result.Results, 1)
		assert.Equal(t, "generation-ii", secondPage.Result.Results[0].Name)
		assert.Equal(t, "https://pokeapi.co/api/v2/generation/2/", secondPage.Result.Results[0].URL)
		assert.Equal(t, 1, secondPage.Result.Count)

		assert.Nil(t, pages.Next(ctx))

		reqs := server.Requests()
		require.Len(t, reqs, 2)
		assert.Equal(t, "/generation", reqs[0].Path)
		assert.Equal(t, map[string][]string{}, reqs[0].Query)
		assert.Equal(t, "/generation", reqs[1].Path)
		assert.Equal(t, map[string][]string{"offset": {"20"}, "limit": {"20"}}, reqs[1].Query)
	})

	t.Run("it paginates through all generations", func(t *testing.T) {
		server.ResetRequests()

		client := server.PokeSDKClient()
		pages := client.Generation.List()

		var pageCount int
		for page := range pages.All(ctx) {
			pageCount++
			require.NotNil(t, page)
			require.Len(t, page.Result.Results, 1)

			switch pageCount {
			case 1:
				assert.Equal(t, "generation-i", page.Result.Results[0].Name)
				assert.Equal(t, "https://pokeapi.co/api/v2/generation/1/", page.Result.Results[0].URL)
			case 2:
				assert.Equal(t, "generation-ii", page.Result.Results[0].Name)
				assert.Equal(t, "https://pokeapi.co/api/v2/generation/2/", page.Result.Results[0].URL)
			default:
				t.Fatalf("unexpected page count: %d", pageCount)
			}
		}
	})
}
