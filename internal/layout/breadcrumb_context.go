package layout

import (
	"encoding/base64"
	"encoding/json"
)

type encodedBreadcrumb struct {
	Title          string `json:"title"`
	URL            string `json:"url,omitempty"`
	IconIdentifier string `json:"iconIdentifier,omitempty"`
}

func EncodeBreadcrumbs(breadcrumbs []Breadcrumb) (string, error) {
	encoded := make([]encodedBreadcrumb, 0, len(breadcrumbs))
	for _, breadcrumb := range breadcrumbs {
		encoded = append(encoded, encodedBreadcrumb{
			Title:          breadcrumb.Title,
			URL:            breadcrumb.URL,
			IconIdentifier: breadcrumb.IconIdentifier,
		})
	}

	payload, err := json.Marshal(encoded)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(payload), nil
}

func DecodeBreadcrumbs(value string) ([]Breadcrumb, error) {
	if value == "" {
		return nil, nil
	}

	payload, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}

	var encoded []encodedBreadcrumb
	if err := json.Unmarshal(payload, &encoded); err != nil {
		return nil, err
	}

	breadcrumbs := make([]Breadcrumb, 0, len(encoded))
	for _, breadcrumb := range encoded {
		breadcrumbs = append(breadcrumbs, Breadcrumb{
			Title:          breadcrumb.Title,
			URL:            breadcrumb.URL,
			IconIdentifier: breadcrumb.IconIdentifier,
		})
	}

	return breadcrumbs, nil
}
