package specification

import "strings"

type ContentType string

const (
	UnknownContentType ContentType = ""
	ApplicationJSON    ContentType = "application/json"
	ApplicationXML     ContentType = "application/xml"
)

func newContentTypeFromString(contentType string) (ContentType, error) {
	switch strings.ToLower(contentType) {
	case "application/json":
		return ApplicationJSON, nil
	case "application/xml":
		return ApplicationXML, nil
	}

	return UnknownContentType, NewNotAllowedContentTypeError(contentType)
}

func (ct ContentType) String() string {
	return string(ct)
}