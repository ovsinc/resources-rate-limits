package errors

import "bytes"

// GetErrorType возвращает тип ошибки. Для НЕ *Error всегда будет "".
func GetErrorType(err error) string {
	if errtype, ok := err.(*Error); ok { //nolint:errorlint
		return errtype.ErrorType().String()
	}

	return ""
}

// GetID возвращает ID ошибки. Для НЕ *Error всегда будет "".
func GetID(err error) (id string) {
	if idtype, ok := err.(*Error); ok { //nolint:errorlint
		return idtype.ID().String()
	}

	return
}

// ErrorOrNil вернет ошибку или nil
// Возможна обработка multierror или одиночной ошибки (*Error, error).
// Если хотя бы одна ошибка в цепочке является ошибкой, то она будет возвращена в качестве результата.
// В противном случае будет возвращен nil.
// Важно: *Error c Severity Warn не является ошибкой.
func ErrorOrNil(err error) error {
	if e, ok := cast(err); ok {
		return e.ErrorOrNil()
	}

	return err
}

func errsFn(errs []error) *Error {
	for _, e := range errs {
		if myerr, ok := simpleCast(e); ok {
			return myerr
		}
	}
	return nil
}

func simpleCast(err error) (*Error, bool) {
	e, ok := err.(*Error) //nolint:errorlint
	return e, ok
}

func cast(err error) (*Error, bool) {
	switch t := err.(type) { //nolint:errorlint
	case interface{ Errors() []error }: // *multiError
		return errsFn(t.Errors()), true

	case interface{ WrappedErrors() []error }: // *github.com/hashicorp/go-multierror.Error
		return errsFn(t.WrappedErrors()), true

	case *Error:
		return t, true
	}

	return nil, false
}

// Cast преобразует тип error в *Error
// Если error не соответствует *Error, то будет создан *Error с сообщением err.Error().
// Для err == nil, вернется nil.
func Cast(err error) *Error {
	if err == nil {
		return nil
	}

	if e, ok := cast(err); ok {
		return e
	}

	return New(err.Error())
}

func findByID(err error, id []byte) (*Error, bool) {
	checkIDFn := func(errs []error) *Error {
		for _, err := range errs {
			if e, ok := simpleCast(err); ok && bytes.Equal(e.ID().Bytes(), id) {
				return e
			}
		}
		return nil
	}

	switch t := err.(type) { //nolint:errorlint
	case interface{ Errors() []error }: // *multiError
		e := checkIDFn(t.Errors())
		return e, e != nil

	case interface{ WrappedErrors() []error }: // *github.com/hashicorp/go-multierror.Error
		e := checkIDFn(t.WrappedErrors())
		return e, e != nil

	case *Error:
		return t, bytes.Equal(t.ID().Bytes(), []byte(id))
	}

	return nil, false
}

// UnwrapByID вернет ошибку (*Error) с указанным ID.
// Если ошибка с указанным ID не найдена, вернется nil.
func UnwrapByID(err error, id string) *Error {
	if e, ok := findByID(err, []byte(id)); ok {
		return e
	}

	return nil
}

// Contains проверит есть ли в цепочке ошибка с указанным ID.
// Допускается в качестве аргумента err указывать одиночную ошибку.
func Contains(err error, id string) bool {
	_, ok := findByID(err, []byte(id))

	return ok
}
