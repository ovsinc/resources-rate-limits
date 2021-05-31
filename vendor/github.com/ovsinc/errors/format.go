package errors

import (
	"fmt"
	"io"
	"sort"
	"strconv"
)

var (
	// DefaultFormatFn функция форматирования, используемая по-умолчанию
	DefaultFormatFn FormatFn //nolint:gochecknoglobals

	// DefaultMultierrFormatFunc функция форматирования для multierr ошибок.
	DefaultMultierrFormatFunc MultierrFormatFn //nolint:gochecknoglobals

	_multilinePrefix    = []byte("the following errors occurred:") //nolint:gochecknoglobals
	_multilineSeparator = []byte("\n")                             //nolint:gochecknoglobals
	_multilineIndent    = []byte("\t#")                            //nolint:gochecknoglobals
	_msgSeparator       = []byte(" -- ")                           //nolint:gochecknoglobals
)

type (
	// FormatFn тип функции форматирования.
	FormatFn func(w io.Writer, e *Error)

	// MultierrFormatFn типу функции морматирования для multierr.
	MultierrFormatFn func(w io.Writer, es []error)
)

// StringFormat функция форматирования вывода сообщения *Error в виде строки.
// Используется по-умолчанию.
func StringFormat(buf io.Writer, e *Error) { //nolint:cyclop
	if e == nil {
		return
	}

	writeDelim := false

	if et := e.ErrorType().Bytes(); len(et) > 0 {
		_, _ = io.WriteString(buf, "(")
		_, _ = buf.Write(et)
		_, _ = io.WriteString(buf, ")")
		writeDelim = true
	}

	if ops := e.Operations(); len(ops) > 0 {
		_, _ = io.WriteString(buf, "[")
		op0 := ops[0]
		_, _ = buf.Write(op0.Bytes())
		for _, opN := range ops[1:] {
			_, _ = io.WriteString(buf, ",")
			_, _ = buf.Write(opN.Bytes())
		}
		_, _ = io.WriteString(buf, "]")
		writeDelim = true
	}

	if ctxs := e.ContextInfo(); len(ctxs) > 0 {
		_, _ = io.WriteString(buf, "<")
		ctxskeys := make([]string, 0, len(ctxs))
		for i := range ctxs {
			ctxskeys = append(ctxskeys, i)
		}
		sort.Strings(ctxskeys)
		_, _ = fmt.Fprintf(buf, "%s:%v", ctxskeys[0], ctxs[ctxskeys[0]])
		for _, i := range ctxskeys[1:] {
			_, _ = io.WriteString(buf, ",")
			_, _ = fmt.Fprintf(buf, "%s:%v", i, ctxs[i])
		}
		_, _ = io.WriteString(buf, ">")
		writeDelim = true
	}

	if writeDelim && len(e.Msg().Bytes()) > 0 {
		_, _ = buf.Write(_msgSeparator)
	}

	_, _ = e.WriteTranslateMsg(buf)
}

//

// multierr

// StringMultierrFormatFunc функция форматирования вывода сообщения для multierr в виде строки.
// Используется по-умолчанию.
func StringMultierrFormatFunc(w io.Writer, es []error) {
	if len(es) == 0 {
		_, _ = io.WriteString(w, "")
		return
	}

	_, _ = w.Write(_multilinePrefix)
	_, _ = w.Write(_multilineSeparator)

	for i, err := range es {
		if err == nil {
			continue
		}
		_, _ = w.Write(_multilineIndent)
		_, _ = io.WriteString(w, strconv.Itoa(i))
		_, _ = io.WriteString(w, " ")
		_, _ = io.WriteString(w, err.Error())
		_, _ = w.Write(_multilineSeparator)
	}
}
