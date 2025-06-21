package encoding

import "net/url"

func EncodeQueryParams(params map[string]string) (string, bool) {
	if len(params) == 0 || params == nil {
		return "", false
	}

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	return "?" + values.Encode(), true
}
