package backend

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/jameshalsall/pokesdk/internal/encoding"
)

const (
	defaultMaxIdleConns        = 100
	defaultTimeout             = 10 * time.Second
	defaultDialTimeout         = 5 * time.Second
	defaultKeepAlive           = 30 * time.Second
	defaultTLSHandshakeTimeout = 5 * time.Second
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTP struct {
	client HTTPClient
}

func NewDefaultHTTP() *HTTP {
	return &HTTP{
		client: defaultHTTPClient(),
	}
}

func NewHTTP(client HTTPClient) *HTTP {
	return &HTTP{
		client: client,
	}
}

// Process sends an HTTP request with the given method and path, marshals params (if present),
// and decodes the response into out (which must be a pointer).
// If a non-pointer is passed as out, an error will be returned.
func (h HTTP) Process(ctx context.Context, url string, params map[string]string, out any) error {

	if params != nil {
		query, ok := encoding.EncodeQueryParams(params)
		if ok {
			url += query
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("pokesdk/backend: failed to create HTTP request: %w", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("pokesdk/backend: HTTP request failed: %w", err)
	}
	defer h.closeResponseBody(resp)

	if resp.StatusCode == http.StatusNotFound {
		return ErrResourceNotFound
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pokesdk/backend: HTTP request failed with status code %d", resp.StatusCode)
	}

	if out != nil {
		if err := encoding.DecodeJSON(resp, out); err != nil {
			return fmt.Errorf("pokesdk/backend: failed to decode HTTP response body: %w", err)
		}
	}

	return nil
}

func (h HTTP) closeResponseBody(resp *http.Response) {
	if resp.Body == nil {
		return
	}

	if err := resp.Body.Close(); err != nil {
		// we should log the error (using a user-provided logger), but the SDK does not support this yet so
		// just spit it out for demonstration purposes...
		fmt.Printf("pokesdk/backend: error closing HTTP response body: %s\n", err)
	}
}

func defaultHTTPClient() HTTPClient {
	return &http.Client{
		Timeout: defaultTimeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   defaultDialTimeout,
				KeepAlive: defaultKeepAlive,
			}).DialContext,
			TLSHandshakeTimeout: defaultTLSHandshakeTimeout,
			MaxIdleConns:        defaultMaxIdleConns,
			MaxIdleConnsPerHost: defaultMaxIdleConns,
		},
	}
}
