package middlewares

import (
	"github.com/ovsinc/multilog"
	rate "github.com/ovsinc/resources-rate-limits"
)

const (
	// HeaderRateLimitRemaining = "X-RateLimit-Remaining"
	HeaderRetryAfter = "Retry-After"
)

type CommonConfig struct {
	// MemoryUsageBarrierPercentage the memory usage barrier after which the rate is enabled
	MemoryUsageBarrierPercentage float64

	// CPUUtilizationBarrierPercentage the CPU utilization barrier after which the rate is enabled
	CPUUtilizationBarrierPercentage float64

	// может быть nil
	Logger multilog.Logger
}

var DefaultCommonConfig = CommonConfig{
	MemoryUsageBarrierPercentage:    rate.DefaultMemoryUsageBarrierPercentage,
	CPUUtilizationBarrierPercentage: rate.DefaultCPUUtilizationBarrierPercentage,
}
