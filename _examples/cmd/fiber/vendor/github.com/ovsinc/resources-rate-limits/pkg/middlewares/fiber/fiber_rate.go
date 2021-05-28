package fiber

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ovsinc/resources-rate-limits/pkg/middlewares"

	rate "github.com/ovsinc/resources-rate-limits"

	sysfiber "github.com/gofiber/fiber/v2"

	"github.com/ovsinc/multilog"
)

const (
	DefaultRetryAfter = "3"
)

type Config struct {
	Rate         rate.Limiter
	Config       *rate.RateLimitConfig
	Logger       multilog.Logger
	ErrHandler   func(*sysfiber.Ctx, *Config) error
	LimitHandler func(*sysfiber.Ctx, *Config, time.Duration) error
	Skip         func(*sysfiber.Ctx) bool
}

var DefaultFiberConfig = Config{
	Config: rate.DefaultRateLimitConfig,
	Skip:   nil,
	// по-дефолту устанавливается *Simple ресорсер и дефолтный конфиг
	Rate: rate.MustNew(),
	// дефолтный ErrHandler установит заголовок "RetryAfter" = DefaultRetryAfter
	ErrHandler: func(ctx *sysfiber.Ctx, _ *Config) error {
		ctx.Response().Header.Add(middlewares.HeaderRetryAfter, DefaultRetryAfter)
		return ctx.SendStatus(http.StatusTooManyRequests)
	},
	// дефолтный LimitHandler установит заголовок
	// "RetryAfter" =  DefaultRetryAfter, если тротлинг (без задержек) или
	// "RetryAfter" = двойное время выполнения *Rate.Limit(), если rate limit
	LimitHandler: func(ctx *sysfiber.Ctx, conf *Config, workingTime time.Duration) error {
		var wtrStr string = DefaultRetryAfter
		wtr := workingTime.Round(time.Second)
		if wtr > 1 {
			wtrStr = strconv.Itoa(2 * int(wtr))
		}
		ctx.Response().Header.Add(middlewares.HeaderRetryAfter, wtrStr)
		return ctx.SendStatus(http.StatusTooManyRequests)
	},
}

// RateLimitWithConfig echo middleware with custom conf
func RateLimitWithConfig(config *Config) sysfiber.Handler {
	if config.Rate == nil {
		config.Rate = rate.MustNew()
	}

	return func(c *sysfiber.Ctx) error {
		if config.Skip != nil && config.Skip(c) {
			return c.Next()
		}

		now := time.Now()

		info := config.Rate.Limit()

		switch {
		case info == nil:
			if config.Logger != nil {
				config.Logger.Errorf("Resource rate limitter fails with no result.")
			}
			return config.ErrHandler(c, config)

		case info.Err != nil:
			if config.Logger != nil {
				config.Logger.Errorf("Resource rate limitter fails with err: '%v'.", info.Err)
			}
			return config.ErrHandler(c, config)

		case info.Time != nil:
			wait := info.Time.Sub(now)

			if config.Logger != nil {
				config.Logger.Warnf(
					"Request from '%v' with path '%v' is rate limited. The request was completed in %s.", c.IP(), c.Path(), wait.String(),
				)
				config.Logger.Infof(
					"Resource rate limite is reached. Memory - %.2f of %.2f, CPU - %.2f of %.2f.",
					info.RAMUsed,
					config.Config.MemoryUsageBarrierPercentage,
					info.CPUUtilization,
					config.Config.CPUUtilizationBarrierPercentage,
				)
			}
			return config.LimitHandler(c, config, wait)
		}

		if config.Logger != nil {
			config.Logger.Debugf(
				"Utilization percents: RAM - %.2f, CPU - %.2f.",
				info.RAMUsed, info.CPUUtilization,
			)
		}

		return c.Next()
	}
}
