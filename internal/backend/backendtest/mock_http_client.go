package backendtest

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.MethodCalled("Do")
	if args.Get(0) != nil {
		return args.Get(0).(*http.Response), nil
	}
	return nil, args.Error(1)
}
