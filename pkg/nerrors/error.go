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

	var target complexError
	if !errors.As(err, &target) {
		return nil
	}

	nestedErrors := target.NestedErrors()
	result := make([]error, len(nestedErrors))
	copy(result, nestedErrors)

	return result
}

func CommonError(err error) string {
	if err == nil {
		return ""
	}

	var target complexError
	if !errors.As(err, &target) {
		return err.Error()
	}

	return target.CommonError()
}
