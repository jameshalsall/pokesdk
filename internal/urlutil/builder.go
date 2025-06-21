package urlutil

import (
	"fmt"
	"strings"
)

func BuildURL(baseURL, path string) string {
	return fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), strings.TrimLeft(path, "/"))
}
