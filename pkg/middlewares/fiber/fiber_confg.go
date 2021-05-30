package fiber

import (
	"net/http"
	"strconv"
	"time"

	sysfiber "github.com/gofiber/fiber/v2"
	rate "github.com/ovsinc/resources-rate-limits"
	"github.com/ovsinc/resources-rate-limits/pkg/middlewares"
)

const (
	DefaultRetryAfter = "3"
	DefaultWaiting    = 3 * time.Second
)

type Config struct {
	middlewares.CommonConfig

	ErrHandler func(*sysfiber.Ctx, *Config, error) error

	LimitHandler func(*sysfiber.Ctx, *Config, time.Time) error

	// может быть nil
	Skip func(*sysfiber.Ctx) bool
}

var DefaultConfig = Config{
	// ErrHandler установит заголовок "RetryAfter" = DefaultRetryAfter
	ErrHandler: func(ctx *sysfiber.Ctx, conf *Config, err error) error {
		if conf.Logger != nil {
			conf.Logger.Errorf(err.Error())
		}

		ctx.Response().Header.Add(middlewares.HeaderRetryAfter, DefaultRetryAfter)
		return ctx.SendStatus(http.StatusTooManyRequests)
	},

	// LimitHandler выполнит замедление запроса на DefaultWaiting
	// и установит заголовок "RetryAfter" = двойное время выполнения *Rate.Limit()
	LimitHandler: func(c *sysfiber.Ctx, conf *Config, t time.Time) error {
		// limit request
		time.Sleep(DefaultWaiting)

		workingTime := time.Since(t)

		var wtrStr string = DefaultRetryAfter
		wtr := workingTime.Round(time.Second)
		if wtr > 1 {
			wtrStr = strconv.Itoa(2 * int(wtr))
		}
		c.Response().Header.Add(middlewares.HeaderRetryAfter, wtrStr)

		if conf.Logger != nil {
			conf.Logger.Warnf(
				"Request from '%v' with path '%v' is rate limited. The request was completed in %s.", c.IP(), c.Path(), workingTime.String(),
			)
		}

		return c.SendStatus(http.StatusTooManyRequests)
	},
}

func defaultRate(limiter rate.Limiter) rate.Limiter {
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
			o.limiter = defaultRate(limiter)
		}
	}
	WithConfig = func(config Config) OptionFiber {
		return func(o *optionFiber) {
			o.config = defaultConfig(config)
		}
	}
)
