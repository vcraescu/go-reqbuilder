package reqbuilder

import (
	"fmt"
	"net/http"
	"net/url"
)

func buildURL(baseURL, path string, rawQuery string) (string, error) {
	rawURL, err := url.JoinPath(baseURL, path)
	if err != nil {
		return "", fmt.Errorf("joinPath: %w", err)
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}

	u.RawQuery = rawQuery

	return u.String(), nil
}

func mergeHeaders(dst http.Header, src ...http.Header) {
	for _, srcHeader := range src {
		for name, values := range srcHeader {
			for _, v := range values {
				dst.Add(name, v)
			}
		}
	}
}

func replaceHeaders(dst http.Header, src ...http.Header) {
	for _, srcHeader := range src {
		for name, values := range srcHeader {
			for _, v := range values {
				dst.Set(name, v)
			}
		}
	}
}
