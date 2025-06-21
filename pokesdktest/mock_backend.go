package pokesdktest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/stretchr/testify/mock"
)

type MockBackend struct {
	sync.Mutex
	mock.Mock

	response []byte
}

func (m *MockBackend) Process(ctx context.Context, path string, params map[string]string, out any) error {
	m.Lock()
	defer m.Unlock()

	if m.response != nil {
		if err := json.NewDecoder(bytes.NewReader(m.response)).Decode(out); err != nil {
			return fmt.Errorf("pokesdktest: error in mock hydrating the response: %w", err)
		}
	}

	return m.MethodCalled("Process", ctx, path, params, out).Error(0)
}

func (m *MockBackend) HydrateWith(response []byte) {
	m.Lock()
	defer m.Unlock()

	m.response = response
}
