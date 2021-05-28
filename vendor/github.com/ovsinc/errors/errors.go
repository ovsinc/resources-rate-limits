package errors

import (
	"fmt"
	"io"

	"github.com/davecgh/go-spew/spew"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/ovsinc/multilog"
	"github.com/valyala/bytebufferpool"
)

var (
	_ error = (*Error)(nil)
	_ interface {
		WithOptions(ops ...Options) *Error

		ID() Objecter
		Severity() Severity
		Msg() Objecter
		ContextInfo() CtxMap
		Operations() Objects
		ErrorType() Objecter

		Sdump() string
		ErrorOrNil() error

		Format(s fmt.State, verb rune)
		Error() string

		TranslateContext() *TranslateContext
		Localizer() *i18n.Localizer
		WriteTranslateMsg(w io.Writer) (int, error)
		TranslateMsg() string

		Log(l ...multilog.Logger)
	} = (*Error)(nil)
)

// Error структура кастомной ошибки.
// Это потоко-безопасный объект.
type Error struct {
	id               Objecter
	msg              Objecter
	severity         Severity
	errorType        Objecter
	operations       Objects
	formatFn         FormatFn
	translateContext *TranslateContext
	localizer        *i18n.Localizer
	contextInfo      CtxMap
}

// New конструктор на необязательных параметрах
// * ops ...Options -- параметризация через функции-парметры.
// См. options.go
//
// ** *Error
func New(msg string, ops ...Options) *Error {
	e := &Error{
		severity: SeverityError,
		msg:      NewObjectFromString(msg),
	}
	for _, op := range ops {
		op(e)
	}
	return e
}

// setters

// WithOptions производит параметризацию *Error с помощью функции-парметры Options.
// Допускается указывать произвольно количество ops.
// Возвращается новый экземпляр *Error.
func (e *Error) WithOptions(ops ...Options) *Error {
	if e == nil {
		return nil
	}

	// copy *Error
	newerr := new(Error)
	*newerr = *e

	for _, op := range ops {
		op(newerr)
	}

	return newerr
}

// getters

// ID возвращает ID ошибки.
// Это безопасный метод, всегда возвращает не nil.
func (e *Error) ID() Objecter {
	if e == nil || e.id == nil {
		return NewObjectEmpty()
	}

	return e.id
}

// Severity возвращает критичность ошибки
func (e *Error) Severity() Severity {
	if e == nil {
		return 0
	}

	return e.severity
}

// Msg возвращает исходное сообщение об ошибке.
// Это безопасный метод, всегда возвращает не nil.
func (e *Error) Msg() Objecter {
	if e == nil || e.msg == nil {
		return NewObjectEmpty()
	}

	return e.msg
}

// ErrorType вернет тип ошибки.
// Это безопасный метод, всегда возвращает не nil.
func (e *Error) ErrorType() Objecter {
	if e == nil || e.errorType == nil {
		return NewObjectEmpty()
	}

	return e.errorType
}

// Operations вернет список операций.
// Это безопасный метод, всегда возвращает не nil.
func (e *Error) Operations() Objects {
	if e == nil || e.operations == nil {
		return NewObjects()
	}
	return e.operations
}

// TranslateContext вернет *TranslateContext.
func (e *Error) TranslateContext() *TranslateContext {
	return e.translateContext
}

// Localizer вернет локализатор *i18n.Localizer.
func (e *Error) Localizer() *i18n.Localizer {
	return e.localizer
}

// Error methods

// Error возвращает строковое представление ошибки.
// Метод для реализации интерфейса error.
// Метод произволит перевод сообщения об ошибки, если localizer != nil.
// Для идентификации сообщения перевода используется ID ошибки.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	var fn FormatFn
	switch {
	case e.formatFn != nil:
		fn = e.formatFn
	case DefaultFormatFn != nil:
		fn = DefaultFormatFn
	default:
		fn = StringFormat
	}

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	fn(buf, e)

	return buf.String()
}

// Format производит форматирование строки, для поддержки fmt.Printf().
func (e *Error) Format(s fmt.State, verb rune) {
	if e == nil {
		return
	}

	switch verb {
	case 'c':
		fmt.Fprintf(s, "%v\n", e.ContextInfo())

	case 'o':
		fmt.Fprintf(s, "%v\n", e.Operations())

	case 'l':
		_, _ = io.WriteString(s, e.Severity().String())

	case 't':
		_, _ = s.Write(e.ErrorType().Bytes())

	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, e.Sdump())
			return
		}
		_, _ = io.WriteString(s, e.Error())

	case 's', 'q':
		_, _ = io.WriteString(s, e.Error())
	}
}

// translate

// WriteTranslateMsg запишет перевод сообщения ошибки в буфер.
// Если не удастся выполнить перевод в буфер w будет записано оригинальное сообщение.
func (e *Error) WriteTranslateMsg(w io.Writer) (int, error) {
	return writeTranslateMsg(e, w)
}

// TranslateMsg вернет перевод сообщения ошибки.
// Если не удастся выполнить перевод, вернет оригинальное сообщение.
func (e *Error) TranslateMsg() string {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	_, _ = e.WriteTranslateMsg(buf)

	return buf.String()
}

// дополнительные методы

// Sdump вернет текстовый дамп ошибки *Error.
func (e *Error) Sdump() string {
	if e == nil {
		return ""
	}

	if e == nil {
		return ""
	}
	return spew.Sdump(e)
}

// ErrorOrNil вернет ошибку или nil.
// ошибкой считается *Error != nil и Severity == SeverityError
// т.е. SeverityWarn НЕ ошибка
func (e *Error) ErrorOrNil() error {
	if e != nil && e.Severity() == SeverityError {
		return e
	}
	return nil
}

// log

// Log выполнит логгирование ошибки с ипользованием логгера l[0].
// Если l не указан, то в качестве логгера будет использоваться логгер по-умолчанию.
func (e *Error) Log(l ...multilog.Logger) {
	customlog(getLogger(l...), e, e.Severity())
}
