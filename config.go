package resourcesratelimits

import (
	"time"

	"github.com/ovsinc/resources-rate-limits/pkg/resources"
)

const (
	DefaultMemoryUsageBarrierPercentage    = 80.0
	DefaultCPUUtilizationBarrierPercentage = 80.0
)

// можно переопределить при сборке
var (
	DefaultRateTimeSeconds    = 3
	DefaultCheckPeroidSeconds = 3
)

type (
	Handler       func(*RateLimitConfig)
	ErrordHandler func(*RateLimitConfig, error)

	// RateLimitConfig rate limit conf
	RateLimitConfig struct {
		// MemoryUsageBarrierPercentage the memory usage barrier after which the rate is enabled
		MemoryUsageBarrierPercentage float64

		// CPUUtilizationBarrierPercentage the CPU utilization barrier after which the rate is enabled
		CPUUtilizationBarrierPercentage float64

		// LimitHandler is called when a request hits the limit
		LimitHandler Handler

		// ErrorHandler is called when a some error happends (memory calculation fails, etc)
		ErrorHandler ErrordHandler
	}
)

var (
	// DefaultRateLimitConfig default config
	DefaultRateLimitConfig = &RateLimitConfig{
		MemoryUsageBarrierPercentage:    DefaultMemoryUsageBarrierPercentage,
		CPUUtilizationBarrierPercentage: DefaultCPUUtilizationBarrierPercentage,
		LimitHandler:                    DefaultLimitedHandler,
		ErrorHandler:                    DefaultThrottleErrorHandler,
	}

	DefaultLimitedHandler = func(cfg *RateLimitConfig) {
		time.Sleep(time.Duration(DefaultRateTimeSeconds) * time.Second)
	}
	DefaultThrottleHandler = func(cfg *RateLimitConfig) {}

	DefaultLimitedErrorHandler = func(cfg *RateLimitConfig, err error) {
		time.Sleep(time.Duration(DefaultRateTimeSeconds) * time.Second)
	}
	DefaultThrottleErrorHandler = func(cfg *RateLimitConfig, err error) {}
)

type Option func(*resourceLimit)

var (
	SetConfig = func(conf *RateLimitConfig) Option {
		return func(rlp *resourceLimit) {
			rlp.conf = conf
		}
	}
	AppendCPUResourcer = func(res resources.Resourcer) Option {
		return func(rlp *resourceLimit) {
			rlp.cpuRes = res
		}
	}
	AppendRAMResourcer = func(res resources.Resourcer) Option {
		return func(rlp *resourceLimit) {
			rlp.ramRes = res
		}
	}
)
