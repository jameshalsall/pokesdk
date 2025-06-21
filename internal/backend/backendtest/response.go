package backendtest

import (
	"io"
	"strings"
)

func NewMockResponseBody(body string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(body))
}
