package specification

import "strings"

type HTTPMethod string

const (
	UnknownHTTPMethod HTTPMethod = ""
	GET               HTTPMethod = "GET"
	POST              HTTPMethod = "POST"
	PUT               HTTPMethod = "PUT"
	PATCH             HTTPMethod = "PATCH"
	DELETE            HTTPMethod = "DELETE"
	OPTIONS           HTTPMethod = "OPTIONS"
	TRACE             HTTPMethod = "TRACE"
	CONNECT           HTTPMethod = "CONNECT"
	HEAD              HTTPMethod = "HEAD"
)

func newHTTPMethodFromString(method string) (HTTPMethod, error) {
	switch strings.ToUpper(method) {
	case "GET":
		return GET, nil
	case "POST":
		return POST, nil
	case "PUT":
		return PUT, nil
	case "PATCH":
		return PATCH, nil
	case "DELETE":
		return DELETE, nil
	case "OPTIONS":
		return OPTIONS, nil
	case "TRACE":
		return TRACE, nil
	case "CONNECT":
		return CONNECT, nil
	case "HEAD":
		return HEAD, nil
	}

	return UnknownHTTPMethod, NewNotAllowedHTTPMethodError(method)
}

func (m HTTPMethod) String() string {
	return string(m)
}
