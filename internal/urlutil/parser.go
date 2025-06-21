package urlutil

import "net/url"

func PathFromURL(rawUrl string) (string, bool) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return "", false
	}
	return parsedUrl.Path + "?" + parsedUrl.RawQuery, true
}
