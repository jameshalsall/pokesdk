package backend

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/jameshalsall/pokesdk/internal/backend/backendtest"
)

func TestHTTP_Process(t *testing.T) {
	ctx := context.Background()

	t.Run("it should process a request and decode it into a response", func(t *testing.T) {
		h, mocks := newHttpForTests(t)

		mocks.client.On("Do", mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       backendtest.NewMockResponseBody(`{"count": 1, "foo": "bar"}`),
		}, nil).Once()

		var response map[string]any
		err := h.Process(ctx, "/foo", map[string]string{"bar": "baz"}, &response)

		assert.NoError(t, err)
		assert.Equal(t, map[string]any{"count": float64(1), "foo": "bar"}, response)
	})

	t.Run("it should process a request when a nil out is provided", func(t *testing.T) {
		h, mocks := newHttpForTests(t)

		mocks.client.On("Do", mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       backendtest.NewMockResponseBody(`{"foo": "bar"}`),
		}, nil).Once()

		err := h.Process(ctx, "/foo", map[string]string{"bar": "baz"}, nil)

		assert.NoError(t, err)
	})

	t.Run("it returns an error if the request fails", func(t *testing.T) {
		h, mocks := newHttpForTests(t)

		mocks.client.On("Do", mock.Anything).Return(nil, assert.AnError).Once()

		var response map[string]any
		err := h.Process(ctx, "/foo", map[string]string{"bar": "baz"}, &response)

		assert.Error(t, err)
		assert.EqualError(t, err, "pokesdk/backend: HTTP request failed: "+assert.AnError.Error())
	})

	t.Run("it returns an error if the response is 404 Not Found", func(t *testing.T) {
		h, mocks := newHttpForTests(t)

		mocks.client.On("Do", mock.Anything).Return(&http.Response{
			StatusCode: http.StatusNotFound,
			Body:       backendtest.NewMockResponseBody(`{"detail": "Not found"}`),
		}, nil).Once()

		var response map[string]any
		err := h.Process(ctx, "/foo", map[string]string{"bar": "baz"}, &response)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrResourceNotFound)
	})

	t.Run("it returns an error if the response cannot be decoded", func(t *testing.T) {
		h, mocks := newHttpForTests(t)

		mocks.client.On("Do", mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       backendtest.NewMockResponseBody(`{{{{`),
		}, nil).Once()

		var response map[string]any
		err := h.Process(ctx, "/foo", map[string]string{"bar": "baz"}, &response)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "pokesdk/backend: failed to decode HTTP response body")
	})

}

type httpMocks struct {
	client *backendtest.MockHTTPClient
}

func newHttpForTests(t *testing.T) (*HTTP, httpMocks) {
	t.Helper()

	mocks := httpMocks{
		client: &backendtest.MockHTTPClient{},
	}
	return NewHTTP(mocks.client), mocks
}
