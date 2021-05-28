package middlewares

const (
	// HeaderRateLimitRemaining = "X-RateLimit-Remaining"
	HeaderRetryAfter = "Retry-After"
)

type Logger interface {
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}
