# Errors

В процессе работы приложения часто приходится возвращать и обрабатывать ошибки. Стандартный пакет ошибок `errors` достаточно беден в плане возможностей. Пакет `github.com/pkg/errors` более интересен, но также не лишен недостатков.

Этот пакет призван добавить возможностей к обработке ошибок. Для удобства использования используется стратегия, принятая в `github.com/pkg/errors` и в целом в [golang](https://golang.org/). Он совместим со стандартным пакетом `errors`.

## Оглавление

1. [Установка](#Установка)
2. [Миграция](#Миграция)
3. [Тестирование](#Тестирование)
   - [Производительность](#Производительность)
4. [Фичи](#Фичи)
5. [Использование](#Использование)
   - [Методы *Error](#Методы-*Error)
   - [Функции-параметры](#Функции-параметры)
   - [Основные хелперы](#Основные-хелперы)
   - [Логгирование ошибки](#Логгирование-ошибки)
   - [Перевод ошибки](#Перевод-ошибки)
   - [Функции форматирования сообщения ошибки](#Функции-форматирования-сообщения-ошибки)
6. [Список задач](#Список-задач)
7. [Лицензия](#Лицензия)

____

## Установка

```text
go get github.com/ovsinc/errors
```

Для простого использования достаточно будет импортировать пакет в своём приложении:

```golang
package main

import (
    "fmt"
    "github.com/ovsinc/errors"
)

func main() {
    fmt.Printf("%v\n", errors.New("hello error"))
}

```

[К оглавлению](#Оглавление)

## Миграция

Поскольку `github.com/ovsinc/errors` совместим с `errors`, то в общем случае миграция достаточно проста.

```golang
package main

import (
    "fmt"
    // "errors"
    "github.com/ovsinc/errors"
)

func main() {
    fmt.Printf("%v\n", errors.New("hello error"))
}

```

[К оглавлению](#Оглавление)

## Тестирование

Склонируйте репозиторий:

```text
git clone https://github.com/ovsinc/errors
cd errors
```

Для запуска юнит-тестов перейдите каталог репозитория в выполните:

```text
make test
```

### Производительность

Для запуска теста производительности перейдите каталог репозитория в выполните:

```text
make bench
```

Для сравнения с аналогичными решениями выполните:

```text
make bench_vendors
```

[К оглавлению](#Оглавление)

## Фичи

- Стандартный интерфейс ошибки (`error`);
- Дополнительные поля описания ошибки: идентификатор (`Objecter`), тип (`Objecter`), контекст (`map[string]interface{}`), severity (Enum: `log.SeveritError`, `log.SeverityWarn`), операции (`Objects`);
- Дополнительный метод: [логгирование](#Логгирование-ошибки);
- Сообщение ошибки может быть [локализовано](#Перевод-ошибки);
- логгирование с помощью [multilog](https://github.com/ovsinc/multilog);
- Логгирование в цепочке;
- Потокобезопасное управления ошибкой/ошибками.

[К оглавлению](#Оглавление)

## Использование

Ошибка `*Error` представлена следующей структурой:

```golang
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
```

Пакет имеет обратную совместимость по методам со стандартной библиотекой `errors` и `github.com/pkg/errors`. Поэтому может быть использован в стандартных сценариях, а также с дополнительными возможностями.

Тип `*Error` является потокобезопасным.

Дополнительные к стандартным могут использованы следующие кейсы:

- оборачивание цепочки ошибок;
- логгирование ошибки в момент ее формирования;
- формирование ошибок в процессе выполнения цепочки методов и проверка в вышестоящем методе (с возможным логгированием);
- выдача ошибок (ошибки) клиентскому приложению с переводом сообщения на язык клиента (при установки локализатора).

Для переводов сообщения используется библиотека `github.com/nicksnyder/go-i18n/v2/i18n`. Ознакомится с особенностями работы i18n можно [тут](https://github.com/nicksnyder/go-i18n).

Можно также ознакомится с [примерами](https://github.com/ovsinc/errors/-/blob/main/example_test.go) использования `github.com/ovsinc/errors`.

### Методы *Error

| Метод | Описание |
| ----- | -------- |
| `func New(msg string, ops ...Options) *Error` | Конструктор `*Error`. Обязательно нужно указать сообщение ошибки. Для ошибки будет установлен `severity = log.SeverityError`. Свойства `*Error` можно установить или переопределить с помощью [функций-параметров](#Функции-параметры). |
| `func NewWithLog(msg string, ops ...Options) *Error` | Конструктор `*Error`, как и `New`. Перед возвратом `*Error` производит логгирование на дефолтном логгере. |
| `func (e *Error) Error() string` | Метод, возвращающий строку сообщения ошибки. |
| `func (e *Error) WithOptions(ops ...Options) *Error` | Вернет новый объект `*Error` как копию `e` с модифицированными с помощью `ops ...Options` свойствами. Изменение свойств `*Error` производится с помощью [функций-параметров](#Функции-параметры). |
| `func (e *Error) Severity() log.Severity` | Геттер получения важности ошибки. |
| `func (e *Error) Msg() Objecter` | Геттер получения объекта сообщения ошибки. |
| `func (e *Error) Sdump() string` | Геттер получения текстового дампа `*Error`. Может использоваться для отладки. Или в выводе форматированного сообщения, например, `fmt.Printf("%+v", err)`. |
| `func (e *Error) ErrorOrNil() error` | Геттер получения ошибки или `nil`. `*Error` с severity `log.SeverityWarn` не является ошибкой; метод `ErrorOrNil` с таким типом вернет `nil`. |
| `func (e *Error) Operations() Objects` | Геттер получения списка объектов операций. |
| `func (e *Error) ErrorType() Objecter` | Геттер получения объекта типа ошибки. |
| `func (e *Error) Format(s fmt.State, verb rune)` | Функция форматирования для обработки строк с возможностью задания формата, например: `fmt.Printf`, `fmt.Sprintf`, `fmt.Fprintf`,.. |
| `func (e *Error) ID() Objecter` | Геттер получения объекта ID. |
| `func (e *Error) TranslateContext() *TranslateContext` | Геттер получения контекста перевода. |
| `func (e *Error) Localizer() *i18n.Localizer` | Геттер получения локализатора. |
| `func (e *Error) Log(l ...multilog.Logger) ` | Метод логгирования. Выполнит логгирование ошибки с использованием логгера `l[0]`. |
| `func (e *Error) WriteTranslateMsg(w io.Writer) (int, error)` | Запишет перевод сообщения ошибки в буфер. В случае неудачи перевода в буфер запишется оригинальное сообщение (без перевода). |
| `func (e *Error) TranslateMsg() string` | Выполнит перевод сообщения и вернет его. В случае неудачи перевода метод вернет оригинальное сообщение (без перевода).|
| `func (e *Error) Is(target error) bool` |  Проверит на равенство `target`. Для go >= 1.13. |
| `func (e *Error) As(target interface{}) bool ` |  Проверит на тип `*Error`. Для go >= 1.13. |
| `func (e *Error) Unwrap() error` | Всегда вернет `nil`. Для go >= 1.13. |

### Объекты (Objects, Objecter)

Тип `Objecter` представлен следующим интерфейсом:

```golang
type Objecter interface {
    String() string
    Bytes() []byte
    Buffer() *bytes.Buffer
}
```

Тип `Objects` можно представить интерфейсом:

```golang
type _ interface {
    Append(oo ...Objecter) Objects
    AppendString(ss ...string) Objects
    AppendBytes(vv ...[]byte) Objects
}
```

В этой версии внесены изменения в методы `ID()`, `Msg()`, `ErrorType()` - изменены тип возвращаемого значения со `string` на `Objecter`.
Например, так: `func (e *Error) ID() string` -> `func (e *Error) ID() Objecter`.

Теперь для получения значения типа `string` достаточно дополнить вызов метода `Objecter.String()`. Это безопасно.

Например:

```golang
func main() {
    e := errors.New("hello")
-   fmt.Println(e.Msg()) 
+   fmt.Println(e.Msg().String()) 
}
```

### Функции-параметры

Параметризация `*Error` производится с помощью функций-параметров типа `type Options func(e *Error)`.

| Метод | Описание |
| ----- | -------- |
| `func SetFormatFn(fn FormatFn) Options` | Устанавливает функцию форматирования. Если значение `nil`, будет использоваться функция форматирования по-умолчанию. |
| `func SetMsg(msg string) Options` | Установить сообщение. |
| `func SetMsgBytes(msg []byte) Options` | Установить сообщение из `[]byte`. |
| `func SetSeverity(severity log.Severity) Options` | Установить уровень важности сообщения. Доступные значения: `log.SeverityWarn`, `log.SeverityError`. |
| `func SetLocalizer(localizer *i18n.Localizer) Options ` | Установить локализатор для перевода. |
| `func SetTranslateContext(tctx *TranslateContext) Options` | Установить `*TranslateContext` для указанного языка. Используется для настройки дополнительных параметров, требуемых для корректного перевода. Например, `TranslateContext.PluralCount` позволяет установить множественное значение используемых в переводе объектов. |
| `func SetErrorType(etype string) Options` | Установить тип ошибки. Тип ошибки - `string`. |
| `func SetErrorTypeBytes(etype []byte) Options` | Установить тип ошибки из `[]byte`. |
| `func SetOperations(ops ...string) Options` | Установить список выполненных операций, указанных как `string`. |
| `func SetOperationsBytes(o ...[]byte) Options` | Установить список выполненных операций, указанных как `[]byte`. |
| `func AppendOperations(ops ...string) Options` | Добавить операции, указанные как `string`, к уже имеющемуся списку. Если список операций не существует, он будет создан. |
| `func AppendOperationsBytes(o ...[]byte) Options` | Добавить операции, указанные как `[]byte`, к уже имеющемуся списку. Если список операций не существует, он будет создан. |
| `func SetContextInfo(ctxinf CtxMap) Options` | Задать контекст ошибки. |
| `func AppendContextInfo(key string, value interface{}) Options` | Добавить значения к уже имеющемуся контексту ошибки. Если контекст ошибки не существует, он будет создан. |
| `func SetID(id string) Options` | Установить ID ошибки. |
| `func SetIDBytes(id []byte) Options` | Установить ID ошибки из `[]byte`. |

### Основные хелперы

Все хелперы работают с типом `error`.

| Хелпер | Описание |
| ------ | -------- |
| `func GetErrorType(err error) string` | Получить тип ошибки. Для НЕ `*Error` всегда будет "". |
| `func ErrorOrNil(err error) error` | Возможна обработка цепочки или одиночной ошибки. Если хотя бы одна ошибка в цепочке является ошибкой, то она будет возвращена в качестве результата. Важно: `*Error` c Severity `Warn` не является ошибкой. |
| `func Cast(err error) *Error` | Преобразование типа `error` в `*Error`. |
| `func Append(errs ...error) error` | Создать цепочку ошибок. Допускается использование `nil` в аргументах. |
| `func Wrap(left error, right error) error` | Обернуть ошибку `left` ошибкой `right`, получив цепочку. Допускается использование `nil` в одном из аргументов, тогда функция вернет ошибку из второго аргумента. |
| `func Errors(err error) []error` | Получить список ошибок из цепочки. Вернет `nil`, при пустой цепочке. |
| `func UnwrapByID(err error, id string) *Error` | Получить ошибку (`*Error`) по ID. Вернет `nil`, если в случае провала поиска. |
| `func GetID(err error) (id string)` | Получить ID ошибки. Для НЕ `*Error` всегда будет "". |
| `func Contains(err error, id string) bool` | Проверить присутствует ли в цепочке ошибка с указанным ID. |
| `func Is(err, target error) bool` | Обёртка над методом стандартной библиотеки `errors.Is`. Для go >= 1.13. |
| `func As(err error, target interface{}) bool` | Обёртка над методом стандартной библиотеки `errors.As`. Для go >= 1.13. |
| `func Unwrap(err error) error` | Обёртка над методом стандартной библиотеки `errors.Unwrap`. Для go >= 1.13. |

### Логгирование ошибки

Логгирование в пакете реализовано с помощью [multilog](https://github.com/ovsinc/multilog).

```golang
type Logger interface {
    Warn(err error)
    Error(err error)
}
```

В пакете [multilog](https://github.com/ovsinc/multilog) присутствует логгер по-умолчанию `https://github.com/ovsinc/multilog.DefaultLogger`.
Он установлен на использование стандартного для Go логгера `log`.

При необходимости его можно легко переопределить на более подходящее значение из пакета [multilog](https://github.com/ovsinc/multilog).

Для логгирования в `*Error` имеется метод `Log(l ...multilog.Logger)`.

Однако, приводить `error` к `*Error` каждый раз не требуется. Для логгирования в пакете есть несколько хелперов.

| Хелпер | Описание |
| ------ | -------- |
| `func NewWithLog(msg string, ops ...Options) *Error` | Функция произведет логгирование ошибки дефолтным логгером. |
| `func Log(err error, l ...multilog.Logger)` | Функция произведет логгирование ошибки дефолтным логгером или логгером указанным в l (будет использоваться только первое значение). |
| `func AppendWithLog(errs ...error) error` | Хелпер создать цепочку ошибок., выполнит логгирование дефолтным логгером и вернет цепочку. |
| `func WrapWithLog(olderr error, err error) error` | Хелпер обернет `olderr` ошибкой `err`, выполнит логгирование дефолтным логгером и вернет цепочку. |

Для удобства поддерживаются несколько оберток над наиболее популярными логгерами.

Ниже приведен пример использования `github.com/ovsinc/errors` c логгированием:

```golang
package main

import (
    "time"

    "github.com/ovsinc/multilog"
    "github.com/ovsinc/multilog/chain"
    "github.com/ovsinc/multilog/journald"
    "github.com/ovsinc/multilog/logrus"
    origlogrus "github.com/sirupsen/logrus"
    "github.com/ovsinc/errors"
)

func main() {
    now := time.Now()

    logrusLogger := logrus.New(origlogrus.New())

    multilog.DefaultLogger = logrusLogger

    err := errors.NewWithLog(
        "hello error",
        errors.SetSeverity(errors.SeverityWarn),
        errors.SetContextInfo(
            errors.CtxMap{
                "time": now,
            },
        ),
    )

    err = err.WithOptions(
        errors.SetID("my id"),
        errors.AppendContextInfo("duration", time.Since(now)),
    )

    journalLogger := journald.New()

    chainLogger := chain.New(logrusLogger, journalLogger)

    err.Log(chainLogger)
}
```

### Перевод ошибки

Для переводов сообщения ошибки используется библиотека `github.com/nicksnyder/go-i18n/v2/i18n`.

Для работы переводов нужно установить:

- `DefaultLocalizer`, тогда он будет использоваться для перевода всех ошибок;
- или локализатор для каждой отдельно взятой ошибки, используя функцию-параметр `*Error.SetLocalizer` при её создании.

Может оказаться удобным установить локализатор `DefaultLocalizer` для всего вашего приложения. Тогда, конечно, ваш локализатор должен содержать весь набор переводимых сообщений и настроен на использование требуемых языков.

В структуре `*Error` за перевод отвечают несколько свойств.

| Свойство | Тип |Назначение | Значение по-умолчанию |
| -------- | --- | --------- | --------------------- |
| translateContext | `*TranslateContext` | Дополнительная информация (контекст) для перевода. | `nil` |
| localizer  | `*i18n.Localizer` | Локализатор. Используется для выполнения переводов сообщения ошибки. | `nil` |

Для выполнения перевода ошибки требуется установить локализатор (если значение `DefaultLocalizer` не было установлено), используя функцию-параметр `SetLocalizer`.
Тогда при вызове метода `*Errors.Error()` будет выдана строка с переведенным сообщением.

В случае возникновения ошибки при переводе сообщения `*Error` будет выдана строка с оригинальным сообщением, без перевода.

Для ошибки `*Error` можно установить контекст перевода. Обычно это требуется для сложных сообщений, например, содержащих имена собственные или количественные значения. Для таких сообщений в составе контекста перевода необходимо установить шаблон `TemplateData map[string]interface{}`.
При использовании множественных форм в сообщении ошибки необходимо установить число в `PluralCount interface{}`.
Можно указать `DefaultMessage *i18n.Message`, если требуется указать значения перевода в случае ошибки перевода из файла.

См. подробности в пакете [i18n](https://github.com/nicksnyder/go-i18n).

Пример использование перевода в сообщении ошибки:

```golang
package main

import (
    _ "embed"
    "fmt"

    "github.com/BurntSushi/toml"
    "github.com/nicksnyder/go-i18n/v2/i18n"
    "github.com/ovsinc/errors"
    "golang.org/x/text/language"
)

//go:embed testdata/active.ru.toml
var translationRu []byte

func main() {
    bundle := i18n.NewBundle(language.English)
    bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
    bundle.MustParseMessageFileBytes(translationRu, "testdata/active.ru.toml")

    err := errors.New(
        "fallback message",
        errors.SetID("ErrEmailsUnreadMsg"),
        errors.SetLocalizer(i18n.NewLocalizer(bundle, "ru")),
        errors.SetTranslateContext(&errors.TranslateContext{
            TemplateData: map[string]interface{}{
                "Name":        "John Snow",
                "PluralCount": 5,
            },
            PluralCount: 5,
        }),
    )

    fmt.Printf("%v\n", err)
}
```

### Функции форматирования сообщения ошибки

В пакете представлены по паре (JSON, string) функций форматирования для единичного сообщения и цепочке сообщений ошибки.

Для цепочки сообщений изменение функции-форматера осуществляется через переменную `DefaultMultierrFormatFunc`. Для неё определено значение по-умолчанию `var DefaultMultierrFormatFunc = StringMultierrFormatFunc`.

Multierror-сообщения форматируются в пакете следующими функциями:

- для вывода в формате JSON - `JSONMultierrFuncFormat(w io.Writer, es []error)`;
- для строкового вывода - `StringMultierrFormatFunc(w io.Writer, es []error)`.

Для сообщений с типом `*Error` используются функции-форматеры типа `type FormatFn func(e *Error) string`. Задать требуемую функцию форматирования можно с помощью функции-параметра `SetFormatFn` в конструкторе или изменить это значение с помощью метода `WithOptions`. Можно задать функцию-форматирования по-умолчанию через переменную `DefaultFormatFn`.

В пакете представлены следующие функции-форматеры:

- для вывода в формате JSON - `JSONFormat(buf io.Writer, e *Error)`;
- для строкового вывода - `StringFormat(buf io.Writer, e *Error)`.

Внимание! При использовании форматирования цепочки сообщения `JSONMultierrFuncFormat` функция форматирование `*Error` по-умолчанию переключается на `JSONFormat`.

Все функций форматирования используют `github.com/valyala/bytebufferpool`, что хорошо сказывается на общей производительности и уменьшает потребление памяти.

[К оглавлению](#Оглавление)

## Список задач

- [ ] Повысить покрытие тестами;
- [ ] Более подробные комментарии для описания методов и функций;
- [ ] Перевод типа ошибки, операций, уровня опасности;
- [ ] Перевод README на en;
- [ ] Выпуск на godoc.

[К оглавлению](#Оглавление)

## Лицензия

Код пакета распространяется под лицензией [Apache 2.0](http://directory.fsf.org/wiki/License:Apache2.0).
