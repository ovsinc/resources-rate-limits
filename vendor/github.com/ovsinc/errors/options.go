package errors

import i18n "github.com/nicksnyder/go-i18n/v2/i18n"

// Options опции из параметра ошибки.
type Options func(e *Error)

// SetFormatFn установит пользовательскую функцию-форматирования.
func SetFormatFn(fn FormatFn) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.formatFn = fn
	}
}

// Msg

// SetMsgBytes установит сообщение об ошибке, указаннов в виде []byte.
func SetMsgBytes(msg []byte) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.msg = NewObjectFromBytes(msg)
	}
}

// SetMsg установит сообщение об ошибке, указанное в виде строки.
func SetMsg(msg string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.msg = NewObjectFromString(msg)
	}
}

// SetSeverity устновит Severity.
func SetSeverity(severity Severity) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.severity = severity
	}
}

func SetSeverityWarn() Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.severity = SeverityWarn
	}
}

func SetSeverityErr() Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.severity = SeverityError
	}
}

// ID

// SetID установит ID ошибки.
func SetID(id string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.id = NewObjectFromString(id)
	}
}

// SetIDBytes установит ID ошибки.
func SetIDBytes(id []byte) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.id = NewObjectFromBytes(id)
	}
}

// Error type

// SetErrorType установит тип ошибки.
func SetErrorType(etype string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.errorType = NewObjectFromString(etype)
	}
}

// SetErrorTypeBytes установит тип ошибки.
func SetErrorTypeBytes(etype []byte) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.errorType = NewObjectFromBytes(etype)
	}
}

// Operations

// AppendOperationsBytes добавить операции, указанные как []byte.
// Можно указать произвольное количество.
// Если в *Error уже были записаны операции,
// то указанные в аргументе будет добавлены к уже имеющимся.
func AppendOperationsBytes(o ...[]byte) Options {
	return func(e *Error) {
		if e == nil || len(o) == 0 {
			return
		}
		e.operations = e.operations.AppendBytes(o...)
	}
}

// AppendOperations добавить операции, указанные как строки.
// Можно указать произвольное количество.
// Если в *Error уже были записаны операции,
// то указанные в аргументе будет добавлены к уже имеющимся.
func AppendOperations(o ...string) Options {
	return func(e *Error) {
		if e == nil || len(o) == 0 {
			return
		}
		e.operations = e.operations.AppendString(o...)
	}
}

// SetOperations установить операции, указанные как строки.
// Можно указать произвольное количество.
// Если в *Error уже были записаны операции,
// то они будут заменены на указанные в аргументе ops.
func SetOperations(o ...string) Options {
	return func(e *Error) {
		if e == nil || len(o) == 0 {
			return
		}
		e.operations = NewObjectsFromString(o...)
	}
}

// SetOperationsBytes установить операции, указанные как []byte.
// Можно указать произвольное количество.
// Если в *Error уже были записаны операции,
// то они будут заменены на указанные в аргументе ops.
func SetOperationsBytes(o ...[]byte) Options {
	return func(e *Error) {
		if e == nil || len(o) == 0 {
			return
		}
		e.operations = NewObjectsFromBytes(o...)
	}
}

// Translate

// SetTranslateContext установит контекст переревода
func SetTranslateContext(tctx *TranslateContext) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.translateContext = tctx
	}
}

// SetLocalizer установит локализатор.
// Этот локализатор будет использован для данной ошибки даже,
// если был установлен DefaultLocalizer.
func SetLocalizer(localizer *i18n.Localizer) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.localizer = localizer
	}
}
