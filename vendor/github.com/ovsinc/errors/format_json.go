package errors

import (
	"fmt"
	"io"
	"strconv"

	json "github.com/goccy/go-json"
)

// JSONFormat функция форматирования вывода сообщения *Error в JSON.
func JSONFormat(buf io.Writer, e *Error) {
	if e == nil {
		_, _ = io.WriteString(buf, "null")
		return
	}

	_, _ = io.WriteString(buf, "{")

	// ID
	_, _ = io.WriteString(buf, "\"id\":")
	_, _ = io.WriteString(buf, "\"")
	_, _ = buf.Write(e.ID().Bytes())
	_, _ = io.WriteString(buf, "\",")

	// ErrorType
	_, _ = io.WriteString(buf, "\"error_type\":")
	_, _ = io.WriteString(buf, "\"")
	_, _ = buf.Write(e.ErrorType().Bytes())
	_, _ = io.WriteString(buf, "\",")

	// Severity
	_, _ = io.WriteString(buf, "\"severity\":")
	_, _ = io.WriteString(buf, "\"")
	_, _ = buf.Write(e.Severity().Bytes())
	_, _ = io.WriteString(buf, "\",")

	// Operations
	_, _ = io.WriteString(buf, "\"operations\":[")
	ops := e.Operations()
	if len(ops) > 0 {
		op0 := ops[0]
		_, _ = io.WriteString(buf, "\"")
		_, _ = buf.Write(op0.Bytes())
		_, _ = io.WriteString(buf, "\"")
		for _, opN := range ops[1:] {
			_, _ = io.WriteString(buf, ",")
			_, _ = io.WriteString(buf, "\"")
			_, _ = buf.Write(opN.Bytes())
			_, _ = io.WriteString(buf, "\"")
		}
	}
	_, _ = io.WriteString(buf, "],")

	// ContextInfo
	_, _ = io.WriteString(buf, "\"context\":")
	cxtInfo := e.ContextInfo()
	if len(cxtInfo) > 0 {
		enc := json.NewEncoder(buf)
		enc.SetIndent("", "")
		_ = enc.Encode(e.ContextInfo())
	} else {
		_, _ = io.WriteString(buf, "null")
	}
	_, _ = io.WriteString(buf, ",")

	// Msg
	_, _ = io.WriteString(buf, "\"msg\":")
	_, _ = io.WriteString(buf, "\"")
	if len(e.Msg().Bytes()) > 0 {
		_, _ = e.WriteTranslateMsg(buf)
	}
	_, _ = io.WriteString(buf, "\"")

	_, _ = io.WriteString(buf, "}")
}

// JSONMultierrFuncFormat функция форматирования вывода сообщения для multierr в виде JSON.
func JSONMultierrFuncFormat(w io.Writer, es []error) {
	if len(es) == 0 {
		_, _ = io.WriteString(w, "null")
	}

	_, _ = io.WriteString(w, "{")

	_, _ = io.WriteString(w, "\"count\":")
	_, _ = io.WriteString(w, strconv.Itoa(len(es)))
	_, _ = io.WriteString(w, ",")

	_, _ = io.WriteString(w, "\"messages\":")
	_, _ = io.WriteString(w, "[")
	writeErrFn := func(e error) {
		if e == nil {
			return
		}
		if myerr, ok := simpleCast(e); ok {
			JSONFormat(w, myerr)
			return
		}
		_, _ = fmt.Fprintf(w, "\"%v\"", e)
	}
	switch len(es) {
	case 1:
		writeErrFn(es[0])
	default:
		writeErrFn(es[0])
		for _, e := range es[1:] {
			_, _ = io.WriteString(w, ",")
			writeErrFn(e)
		}
	}
	_, _ = io.WriteString(w, "]")

	_, _ = io.WriteString(w, "}")
}
