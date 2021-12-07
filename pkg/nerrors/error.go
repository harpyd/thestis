package nerrors

import "github.com/pkg/errors"

type complexError interface {
	CommonError() string
	NestedErrors() []error
}

func NestedErrors(err error) []error {
	if err == nil {
		return nil
	}

	var cerr complexError
	if !errors.As(err, &cerr) {
		return nil
	}

	nestedErrors := cerr.NestedErrors()
	result := make([]error, len(nestedErrors))
	copy(result, nestedErrors)

	return result
}

func CommonError(err error) string {
	if err == nil {
		return ""
	}

	var cerr complexError
	if !errors.As(err, &cerr) {
		return err.Error()
	}

	return cerr.CommonError()
}
