package specification

import "strings"

type HTTPMethod string

const (
	UnknownHTTPMethod HTTPMethod = "!"
	EmptyHTTPMethod   HTTPMethod = ""
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
	methods := map[string]HTTPMethod{
		EmptyHTTPMethod.String(): EmptyHTTPMethod,
		GET.String():             GET,
		POST.String():            POST,
		PUT.String():             PUT,
		PATCH.String():           PATCH,
		DELETE.String():          DELETE,
		OPTIONS.String():         OPTIONS,
		TRACE.String():           TRACE,
		CONNECT.String():         CONNECT,
		HEAD.String():            HEAD,
	}

	if m, ok := methods[strings.ToUpper(method)]; ok {
		return m, nil
	}

	return UnknownHTTPMethod, NewNotAllowedHTTPMethodError(method)
}

func (m HTTPMethod) String() string {
	return string(m)
}
