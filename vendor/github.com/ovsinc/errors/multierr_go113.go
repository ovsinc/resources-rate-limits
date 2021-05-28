// +build go1.13

package errors

func (merr *multiError) Unwrap() error {
	if merr.count < 1 || merr.curIdx > merr.count {
		return nil
	}

	err := merr.errors[merr.curIdx]
	merr.curIdx++

	return err
}

func (merr *multiError) As(target interface{}) bool {
	if x, ok := target.(*Multierror); ok { //nolint:errorlint
		*x = merr
		return true
	}
	return false
}

func (merr *multiError) Is(target error) bool {
	if x, ok := target.(Multierror); ok { //nolint:errorlint
		return x == merr
	}
	return false
}
