//go:build integration

package integration

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/jameshalsall/pokesdk"
)

// Response represents a stubbed response for a given path.
type Response struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
}

// RequestCapture captures the incoming HTTP GET request.
type RequestCapture struct {
	Path    string
	Headers http.Header
	Query   map[string][]string
}

// MockServer is a simple mock HTTP server for GET requests.
type MockServer struct {
	t         *testing.T
	server    *httptest.Server
	mu        sync.Mutex
	responses map[string]Response
	captured  []RequestCapture
}

// NewMockServer initializes and starts the test server.
func NewMockServer(t *testing.T) *MockServer {
	ms := &MockServer{
		t:         t,
		responses: make(map[string]Response),
	}

	ms.server = httptest.NewServer(http.HandlerFunc(ms.handler))
	return ms
}

func (ms *MockServer) PokeSDKClient() *pokesdk.Client {
	return pokesdk.NewClient(pokesdk.WithCustomBaseURL(ms.URL()), pokesdk.WithCustomHttpClient(ms.server.Client()))
}

// Close shuts down the server.
func (ms *MockServer) Close() {
	ms.server.Close()
}

// URL returns the base URL of the test server.
func (ms *MockServer) URL() string {
	return ms.server.URL
}

// StubGET registers a mock response for a given path.
func (ms *MockServer) StubGET(path string, resp Response) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.responses[path] = resp
}

// Requests returns a copy of all captured GET requests.
func (ms *MockServer) Requests() []RequestCapture {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	out := make([]RequestCapture, len(ms.captured))
	copy(out, ms.captured)

	return out
}

func (ms *MockServer) Reset() {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.responses = make(map[string]Response)
	ms.captured = nil
}

func (ms *MockServer) ResetRequests() {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.captured = nil
}

func (ms *MockServer) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET supported", http.StatusMethodNotAllowed)
		return
	}

	ms.mu.Lock()
	url := r.URL.Path
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}
	resp, ok := ms.responses[url]
	ms.captured = append(ms.captured, RequestCapture{
		Path:    r.URL.Path,
		Headers: r.Header.Clone(),
		Query:   r.URL.Query(),
	})
	ms.mu.Unlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	for k, v := range resp.Headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(resp.Body)
}
