package echo

import (
	"net/http"
	"strconv"
	"time"

	sysecho "github.com/labstack/echo/v4"
	rate "github.com/ovsinc/resources-rate-limits"
	"github.com/ovsinc/resources-rate-limits/pkg/middlewares"
)

type Config struct {
	middlewares.CommonConfig

	ErrHandler func(sysecho.Context, *Config, error) error

	LimitHandler func(sysecho.Context, *Config, time.Time) error

	// может быть nil
	Skip func(sysecho.Context) bool
}

var DefaultConfig = Config{
	// ErrHandler установит заголовок "RetryAfter" = DefaultRetryAfter
	ErrHandler: func(ctx sysecho.Context, conf *Config, err error) error {
		middlewares.ErrorThrottleHandler(
			conf.Logger,
			middlewares.Client{
				IP:   ctx.RealIP(),
				Path: ctx.Path(),
			},
			err,
			time.Now(),
		)
		ctx.Response().Header().Add(middlewares.HeaderRetryAfter, middlewares.DefaultRetryAfter)
		return ctx.NoContent(http.StatusTooManyRequests)
	},

	// LimitHandler выполнит замедление запроса на DefaultWaiting
	// и установит заголовок "RetryAfter" = двойное время выполнения *Rate.Limit()
	LimitHandler: func(c sysecho.Context, conf *Config, now time.Time) error {
		middlewares.LimitHandler(
			conf.Logger,
			middlewares.Client{
				IP:   c.RealIP(),
				Path: c.Path(),
			},
			now,
		)

		var wtrStr string = middlewares.DefaultRetryAfter
		wtr := time.Since(now).Round(time.Second)
		if wtr > 1 {
			wtrStr = strconv.Itoa(2 * int(wtr))
		}
		c.Response().Header().Add(middlewares.HeaderRetryAfter, wtrStr)

		return c.NoContent(http.StatusTooManyRequests)
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
