package fiber

import (
	"net/http"
	"strconv"
	"time"

	sysfiber "github.com/gofiber/fiber/v2"
	rate "github.com/ovsinc/resources-rate-limits"
	"github.com/ovsinc/resources-rate-limits/pkg/middlewares"
)

type Config struct {
	middlewares.CommonConfig

	ErrHandler func(*sysfiber.Ctx, *Config, error) error

	LimitHandler func(*sysfiber.Ctx, *Config, time.Time) error

	// может быть nil
	Skip func(*sysfiber.Ctx) bool
}

var DefaultConfig = Config{
	CommonConfig: middlewares.DefaultCommonConfig,

	// ErrHandler установит заголовок "RetryAfter" = DefaultRetryAfter
	ErrHandler: func(ctx *sysfiber.Ctx, conf *Config, err error) error {
		middlewares.ErrorThrottleHandler(
			conf.Logger,
			middlewares.Client{
				IP:   ctx.IP(),
				Path: ctx.Path(),
			},
			err,
			time.Now(),
		)
		ctx.Response().Header.Add(middlewares.HeaderRetryAfter, middlewares.DefaultRetryAfter)
		return ctx.SendStatus(http.StatusTooManyRequests)
	},

	// LimitHandler выполнит замедление запроса на DefaultWaiting
	// и установит заголовок "RetryAfter" = двойное время выполнения *Rate.Limit()
	LimitHandler: func(c *sysfiber.Ctx, conf *Config, now time.Time) error {
		middlewares.LimitHandler(
			conf.Logger,
			middlewares.Client{
				IP:   c.IP(),
				Path: c.Path(),
			},
			now,
		)

		var wtrStr string = middlewares.DefaultRetryAfter
		wtr := time.Since(now).Round(time.Second)
		if wtr > time.Second {
			wtrStr = strconv.Itoa(2 * int(wtr.Seconds()))
		}
		c.Response().Header.Set(middlewares.HeaderRetryAfter, wtrStr)

		return c.SendStatus(http.StatusTooManyRequests)
	},
}

func defaultLimiter(limiter rate.Limiter) rate.Limiter {
	if limiter == nil {
		return rate.MustNew()
	}
	return limiter
}

func defaultConfig(cfg Config) Config {
	if cfg.MemoryUsageBarrierPercentage == 0 {
		cfg.MemoryUsageBarrierPercentage = DefaultConfig.MemoryUsageBarrierPercentage
	}

	if cfg.CPUUtilizationBarrierPercentage == 0 {
		cfg.CPUUtilizationBarrierPercentage = DefaultConfig.CPUUtilizationBarrierPercentage
	}

	if cfg.ErrHandler == nil {
		cfg.ErrHandler = DefaultConfig.ErrHandler
	}

	if cfg.LimitHandler == nil {
		cfg.LimitHandler = DefaultConfig.LimitHandler
	}

	// Skip = nil
	// Logger = nil

	return cfg
}

type optionFiber struct {
	config  Config
	limiter rate.Limiter
}

type OptionFiber func(*optionFiber)

var (
	WithLimiter = func(limiter rate.Limiter) OptionFiber {
		return func(o *optionFiber) {
			o.limiter = limiter
		}
	}
	WithConfig = func(config Config) OptionFiber {
		return func(o *optionFiber) {
			o.config = config
		}
	}
)
