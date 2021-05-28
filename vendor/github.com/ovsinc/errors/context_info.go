package errors

// CtxMap map контекста ошибки.
// В качестве ключа всегда должна быть строка, а значение - любой тип.
// При преобразовании ошибки в строку CtxMap может использоваться различные методы.
// Для функции JSONFormat CtxMap будет преобразовываться с помощью JSON marshall.
// Для функции StringFormat CtxMap будет преобразовываться с помощью fmt.Sprintf.
type CtxMap map[string]interface{}

// SetContextInfo установить CtxMap.
func SetContextInfo(ctxinf CtxMap) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.contextInfo = ctxinf
	}
}

// AppendContextInfo добавить в имеющийся CtxMap значение value по ключу key.
// Если CtxMap в *Error не установлен, то он будет предварительно установлен.
func AppendContextInfo(key string, value interface{}) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		if e.contextInfo == nil {
			e.contextInfo = make(CtxMap)
		}
		e.contextInfo[key] = value
	}
}

//

// ContextInfo вернет контекст CtxMap ошибки.
func (e *Error) ContextInfo() CtxMap {
	return e.contextInfo
}
