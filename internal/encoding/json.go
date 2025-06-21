package encoding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func DecodeJSON(resp *http.Response, into any) error {
	if err := json.NewDecoder(resp.Body).Decode(into); err != nil {
		return fmt.Errorf("pokesdk: failed to decode JSON for type %T: %w", into, err)
	}
	return nil
}
func EncodeJSON(from any) (*bytes.Reader, error) {
	data, err := json.Marshal(from)
	if err != nil {
		return bytes.NewReader(nil), fmt.Errorf("pokesdk: failed to encode JSON for type %T: %w", from, err)
	}
	return bytes.NewReader(data), nil
}
