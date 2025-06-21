package encoding

import (
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
