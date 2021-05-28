package errors

import (
	"runtime"
	"strconv"
	"strings"
)

//
// From https://github.com/go-kit/kit/blob/master/log/value.go
//

// Caller returns a Valuer that returns a file and line from a specified depth
// in the callstack. Users will probably want to use DefaultCaller.
func Caller(depth int) func() string {
	return func() string {
		_, file, line, _ := runtime.Caller(depth)
		idx := strings.LastIndexByte(file, '/')
		// using idx+1 below handles both of following cases:
		// idx == -1 because no "/" was found, or
		// idx >= 0 and we want to start at the character after the found "/".
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}

// DefaultCaller is a Valuer that returns the file and line where the Log
// method was invoked. It can only be used with log.With.
var DefaultCaller = Caller(3) //nolint:gochecknoglobals
