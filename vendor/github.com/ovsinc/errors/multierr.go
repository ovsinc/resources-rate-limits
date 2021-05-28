// Используется оригинальный код проекта "go.uber.org/multierr" с частичным заимствованием.
// Код проекта "go.uber.org/multierr" распространяется под лицензией MIT (https://github.com/uber-go/multierr/blob/master/LICENSE.txt).

package errors

import (
	"fmt"
	"io"

	"github.com/valyala/bytebufferpool"
)

// Errors returns a slice containing zero or more errors that the supplied
// error is composed of. If the error is nil, a nil slice is returned.
//
// 	err := multierr.Append(r.Close(), w.Close())
// 	errors := multierr.Errors(err)
//
// If the error is not composed of other errors, the returned slice contains
// just the error that was passed in.
//
// Callers of this function are free to modify the returned slice.
func Errors(err error) []error {
	if err == nil {
		return nil
	}

	// Note that we're casting to multiError, not errorGroup. Our contract is
	// that returned errors MAY implement errorGroup. Errors, however, only
	// has special behavior for multierr-specific error objects.
	//
	// This behavior can be expanded in the future but I think it's prudent to
	// start with as little as possible in terms of contract and possibility
	// of misuse.
	if eg, ok := err.(*multiError); ok { //nolint:errorlint
		errors := eg.Errors()
		result := make([]error, len(errors))
		copy(result, errors)
		return result
	}

	return []error{err}
}

var _ Multierror = (*multiError)(nil)

type Multierror interface {
	Errors() []error
	Error() string
	Format(f fmt.State, c rune)
}

// multiError is an error that holds one or more errors.
//
// An instance of this is guaranteed to be non-empty and flattened. That is,
// none of the errors inside multiError are other multiErrors.
//
// multiError formats to a semi-colon delimited list of error messages with
// %v and with a more readable multi-line format with %+v.
type multiError struct {
	curIdx int
	count  int
	errors []error
}

// Errors returns the list of underlying errors.
//
// This slice MUST NOT be modified.
func (merr *multiError) Errors() []error {
	if merr == nil {
		return nil
	}
	return merr.errors
}

func (merr *multiError) Error() string {
	if merr == nil {
		return ""
	}

	buff := bytebufferpool.Get()
	defer bytebufferpool.Put(buff)

	merr.writeLines(buff)

	return buff.String()
}

func (merr *multiError) Format(f fmt.State, c rune) {
	switch c {
	case 'w', 'v', 's':
		merr.writeLines(f)
	case 'j':
		JSONMultierrFuncFormat(f, merr.errors)
	}
}

func (merr *multiError) writeLines(w io.Writer) {
	if DefaultMultierrFormatFunc == nil {
		StringMultierrFormatFunc(w, merr.errors)
		return
	}
	DefaultMultierrFormatFunc(w, merr.errors)
}

type inspectResult struct {
	// Number of top-level non-nil errors
	Count int

	// Total number of errors including multiErrors
	Capacity int
}

// Inspects the given slice of errors so that we can efficiently allocate
// space for it.
func inspect(errors []error) (res inspectResult) {
	for _, err := range errors {
		if err == nil {
			continue
		}

		res.Count++
		if merr, ok := err.(*multiError); ok { //nolint:errorlint
			res.Capacity += len(merr.errors)
		} else {
			res.Capacity++
		}
	}
	return
}

// fromSlice converts the given list of errors into a single error.
func fromSlice(errors []error) error {
	res := inspect(errors)
	if res.Count == 0 {
		return nil
	}

	nonNilErrs := make([]error, 0, res.Capacity)
	for _, err := range errors {
		if err == nil {
			continue
		}

		if nested, ok := err.(*multiError); ok { //nolint:errorlint
			nonNilErrs = append(nonNilErrs, nested.errors...)
		} else {
			nonNilErrs = append(nonNilErrs, err)
		}
	}

	return &multiError{errors: nonNilErrs, count: len(nonNilErrs)}
}

// Append создаст цепочку ошибок из ошибок ...errors.
// Допускается использование `nil` в аргументах.
func Append(errors ...error) error {
	return fromSlice(errors)
}

// Wrap обернет ошибку `left` ошибкой `right`, получив цепочку.
// Допускается использование `nil` в одном из аргументов.
func Wrap(left error, right error) error {
	switch {
	case left == nil:
		return right
	case right == nil:
		return left
	}
	return fromSlice([]error{left, right})
}
