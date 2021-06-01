package fiber

import (
	sysfiber "github.com/gofiber/fiber/v2"
	rate "github.com/ovsinc/resources-rate-limits"
	"github.com/ovsinc/resources-rate-limits/pkg/middlewares"
)

// RateLimitWithConfig echo middleware with custom conf
func RateLimit(ops ...OptionFiber) sysfiber.Handler {
	op := new(optionFiber)
	for _, f := range ops {
		f(op)
	}

	// установим дефолтные значения, если они не были заданы
	op.config = defaultConfig(op.config)
	op.limiter = defaultLimiter(op.limiter)
	// настроим дебаг
	if op.config.Debug {
		op.limiter = op.limiter.With(
			rate.SetDebug(op.config.Debug),
		)
	}

	return func(c *sysfiber.Ctx) error {
		if op.config.Skip != nil && op.config.Skip(c) {
			return c.Next()
		}

		info := op.limiter.Limit()

		switch {
		case info == nil:
			return op.config.ErrHandler(c, &op.config, middlewares.ErrResourcerNoResult)

		case info.CPUUtilization == rate.FailValue:
			return op.config.ErrHandler(c, &op.config, middlewares.ErrGetCPUUtilizationFail)

		case info.CPUUtilization == rate.DoneValue || info.RAMUsed == rate.DoneValue:
			return op.config.ErrHandler(c, &op.config, middlewares.ErrClosedResourcer)

		case info.RAMUsed == rate.FailValue:
			return op.config.ErrHandler(c, &op.config, middlewares.ErrGetRAMUsedFail)

		case info.RAMUsed >= op.config.MemoryUsageBarrierPercentage,
			info.CPUUtilization >= op.config.CPUUtilizationBarrierPercentage:
			if op.config.Logger != nil {
				op.config.Logger.Infof(
					"Resource rate limite is reached. Memory - %.2f of %.2f, CPU - %.2f of %.2f.",
					info.RAMUsed,
					op.config.MemoryUsageBarrierPercentage,
					info.CPUUtilization,
					op.config.CPUUtilizationBarrierPercentage,
				)
			}
			if op.config.LimitHandler != nil {
				return op.config.LimitHandler(c, &op.config, info.Time)
			}
		}

		if op.config.Logger != nil {
			if op.config.Debug {
				op.config.Logger.Debugf(
					"Utilization percents: RAM - %.2f, CPU - %.2f.",
					info.RAMUsed, info.CPUUtilization,
				)
			}
		}

		return c.Next()
	}
}
