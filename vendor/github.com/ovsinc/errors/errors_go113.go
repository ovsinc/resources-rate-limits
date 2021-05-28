// +build go1.13

package errors

var _ errorergo13 = (*Error)(nil)

type errorergo13 interface {
	Is(target error) bool
	As(target interface{}) bool
	Unwrap() error
}

func (e *Error) Is(target error) bool {
	switch x := target.(type) { //nolint:errorlint
	case *Error:
		return e == x
	}

	return false
}

func (e *Error) As(target interface{}) bool {
	switch x := target.(type) { //nolint:errorlint
	case **Error:
		*x = e

	default:
		return false
	}

	return true
}

func (e *Error) Unwrap() error {
	return nil
}
