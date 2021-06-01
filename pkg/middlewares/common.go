package middlewares

import (
	"time"

	"github.com/ovsinc/multilog"
	rate "github.com/ovsinc/resources-rate-limits"
)

const (
	// HeaderRateLimitRemaining = "X-RateLimit-Remaining"
	HeaderRetryAfter = "Retry-After"
)

const (
	DefaultRetryAfter = "3"
	DefaultWaiting    = 3 * time.Second
)

type CommonConfig struct {
	// MemoryUsageBarrierPercentage the memory usage barrier after which the rate is enabled
	MemoryUsageBarrierPercentage float64

	// CPUUtilizationBarrierPercentage the CPU utilization barrier after which the rate is enabled
	CPUUtilizationBarrierPercentage float64

	// может быть nil
	Logger multilog.Logger

	Debug bool
}

var DefaultCommonConfig = CommonConfig{
	MemoryUsageBarrierPercentage:    rate.DefaultMemoryUsageBarrierPercentage,
	CPUUtilizationBarrierPercentage: rate.DefaultCPUUtilizationBarrierPercentage,
	Debug:                           false,
}

type Client struct {
	IP   string
	Path string
}

func LimitHandler(logger multilog.Logger, c Client, now time.Time) {
	// limit request
	time.Sleep(DefaultWaiting)
	if logger != nil {
		logger.Warnf(
			"Limited. Request from '%v' with path '%v' is rate limited. The request was completed in %s.",
			c.IP, c.Path, time.Since(now).String(),
		)
	}
}

func ThrottleHandler(logger multilog.Logger, c Client, now time.Time) {
	if logger != nil {
		logger.Warnf(
			"Throttled. Request from '%v' with path '%v' is throttled. The request was completed in %s.",
			c.IP, c.Path, time.Since(now).String(),
		)
	}
}

func ErrorThrottleHandler(logger multilog.Logger, c Client, err error, now time.Time) {
	if logger != nil {
		logger.Errorf(
			"Throttled. Error. Request from '%v' with path '%v' was fails with '%v'.  The request was completed in %s.",
			c.IP, c.Path, err.Error(), time.Since(now).String(),
		)
	}
}

func ErrorLimitHandler(logger multilog.Logger, c Client, err error, now time.Time) {
	// limit request
	time.Sleep(DefaultWaiting)
	if logger != nil {
		logger.Errorf(
			"Limited. Error. Request from '%v' with path '%v' was fails with '%v'. The request was completed in %s.",
			c.IP, c.Path, err.Error(), time.Since(now).String(),
		)
	}
}
