package errors

import (
	"io"

	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

// DefaultLocalizer локализатор по-умолчанию.
// Для каждой ошибки можно переопределить локализатор.
var DefaultLocalizer *i18n.Localizer //nolint:gochecknoglobals

// TranslateContext контекст перевода. Не является обязательным для корректного перевода.
type TranslateContext struct {
	// TemplateData - map для замены в шаблоне
	TemplateData map[string]interface{}
	// PluralCount признак множественности.
	// Может иметь значение nil или число.
	PluralCount interface{}
	// DefaultMessage сообщение, которое будет использовано при ошибке перевода.
	DefaultMessage *i18n.Message
}

func writeTranslateMsg(e *Error, w io.Writer) (int, error) {
	buf := e.Msg().Bytes()

	if len(buf) == 0 {
		return 0, nil
	}

	var localizer *i18n.Localizer
	switch {
	case e.localizer != nil:
		localizer = e.localizer
	case DefaultLocalizer != nil:
		localizer = DefaultLocalizer
	}

	if localizer == nil {
		return w.Write(buf)
	}

	i18nConf := &i18n.LocalizeConfig{
		MessageID: e.id.String(),
	}
	if e.translateContext != nil {
		i18nConf.DefaultMessage = e.translateContext.DefaultMessage
		i18nConf.PluralCount = e.translateContext.PluralCount
		i18nConf.TemplateData = e.translateContext.TemplateData
	}

	msg, _, err := e.localizer.LocalizeWithTag(i18nConf)
	if err != nil {
		return w.Write(buf)
	}

	return io.WriteString(w, msg)
}
