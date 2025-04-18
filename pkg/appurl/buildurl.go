package appurl

import "net/url"

// BuildURL builds a safe URL with optional query parameters
func BuildURL(basePath string, pathParams []string, queryParams map[string]string) string {
	u := &url.URL{
		Path: basePath,
	}

	// Append path parameters if any
	for _, param := range pathParams {
		u.Path += "/" + url.PathEscape(param)
	}

	// Add query parameters
	q := url.Values{}
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return u.String()
}
