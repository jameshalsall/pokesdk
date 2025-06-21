package pokesdk

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPage struct {
	id   int
	next string
}

func (m *mockPage) GetNextURL() string {
	return m.next
}

func TestPaginator_Next(t *testing.T) {
	t.Run("it gets a single page", func(t *testing.T) {
		var fetchCount int
		fetch := func(ctx context.Context, url string) (*mockPage, error) {
			fetchCount++
			return &mockPage{id: 1, next: ""}, nil
		}

		p := NewPaginator[*mockPage]("start", fetch)

		page := p.Next(context.Background())
		require.NoError(t, page.Error)
		assert.NotNil(t, page)
		assert.Equal(t, 1, page.Result.id)

		page = p.Next(context.Background())
		assert.Nil(t, page)

		assert.Equal(t, 1, fetchCount)
	})

	t.Run("it returns an error from pagination", func(t *testing.T) {
		expectedErr := errors.New("fetch failed")

		fetch := func(ctx context.Context, url string) (*mockPage, error) {
			return nil, expectedErr
		}

		p := NewPaginator[*mockPage]("start", fetch)

		page := p.Next(context.Background())
		require.Error(t, page.Error)
		assert.EqualError(t, page.Error, expectedErr.Error())
	})
}

func TestPaginator_All(t *testing.T) {
	t.Run("it paginates multiple pages", func(t *testing.T) {
		pagesData := []*mockPage{
			{id: 1, next: "url2"},
			{id: 2, next: "url3"},
			{id: 3, next: ""},
		}

		var index int
		fetch := func(ctx context.Context, url string) (*mockPage, error) {
			if index >= len(pagesData) {
				return nil, nil
			}
			page := pagesData[index]
			index++
			return page, nil
		}

		p := NewPaginator[*mockPage]("start", fetch)
		ctx := context.Background()

		ch := p.All(ctx)

		var results []*mockPage
		for page := range ch {
			assert.NoError(t, page.Error)
			results = append(results, page.Result)
		}

		assert.Len(t, results, 3)
	})

	t.Run("it does nothing if the context cancelled", func(t *testing.T) {
		fetch := func(ctx context.Context, url string) (*mockPage, error) {
			time.Sleep(50 * time.Millisecond)
			return &mockPage{id: 1, next: ""}, nil
		}

		p := NewPaginator[*mockPage]("start", fetch)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		ch := p.All(ctx)

		select {
		case _, ok := <-ch:
			assert.False(t, ok, "expected pages channel to be closed")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("timeout waiting for pages channel to close")
		}
	})
}
