package fiber_test

import (
	"net/http"
	"strconv"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	rate "gitlab.com/ovsinc/memory-rate-limits"
	"gitlab.com/ovsinc/memory-rate-limits/middlewares"
	middlfiber "gitlab.com/ovsinc/memory-rate-limits/middlewares/fiber"
	memmock "gitlab.com/ovsinc/memory-rate-limits/resources/memory/mock"
)

func newApp(h fiber.Handler) (*fiber.App, *fasthttp.RequestCtx) {
	app := fiber.New()

	if h != nil {
		app.Use(h)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello tester!")
	})

	fctx := &fasthttp.RequestCtx{}
	fctx.Request.Header.SetMethod("GET")
	fctx.Request.SetRequestURI("/")

	return app, fctx
}

func TestRateLimitDefault(t *testing.T) {
	t.Parallel()

	app, fctx := newApp(middlfiber.RateLimit())

	assert := assert.New(t)

	app.Handler()(fctx)

	assert.Equal(http.StatusOK, fctx.Response.StatusCode())
}

func TestRateLimitWithConfigLimited(t *testing.T) {
	t.Parallel()

	limCfg := middlfiber.DefaultFiberConfig
	throttleSecond := 3
	limCfg.Rate = rate.New(
		rate.WithConfig(
			&rate.RateLimitConfig{
				MemoryUsageBarrierPercentage: 80,
				RateSpeed:                    throttleSecond,
				RetryAfter:                   rate.DefailtRetryAfter,
				LimitReached:                 rate.DefaultLimitedHandler,
				ErrorHandler:                 rate.DefaultThrottleErrorHandler,
				Limiter:                      rate.DefaultLimiter,
			},
		),
		rate.WithMemoryChecker(memmock.NewMemLimited()),
	)

	app, fctx := newApp(middlfiber.RateLimitWithConfig(&limCfg))

	assert := assert.New(t)

	app.Handler()(fctx)

	assert.Equal(http.StatusTooManyRequests, fctx.Response.StatusCode())
	assert.Equal(
		strconv.Itoa(throttleSecond),
		string(fctx.Response.Header.Peek(middlewares.HeaderRetryAfter)),
	)
}

func TestRateLimitWithConfigFail(t *testing.T) {
	t.Parallel()

	failCfg := middlfiber.DefaultFiberConfig
	failCfg.Rate = rate.New(
		rate.WithMemoryChecker(memmock.NewMemFails()),
	)
	failCfg.ErrHandler = func(ctx *fiber.Ctx, _ *middlfiber.Config) error {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	app, fctx := newApp(middlfiber.RateLimitWithConfig(&failCfg))

	assert := assert.New(t)

	app.Handler()(fctx)

	assert.Equal(http.StatusInternalServerError, fctx.Response.StatusCode())
}
