package fiber_test

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	fiblogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ovsinc/multilog/golog"
	rate "github.com/ovsinc/resources-rate-limits"
	"github.com/ovsinc/resources-rate-limits/pkg/middlewares"
	middlfiber "github.com/ovsinc/resources-rate-limits/pkg/middlewares/fiber"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
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

func TestRateLimitDefaultOK(t *testing.T) {
	t.Parallel()

	app, fctx := newApp(middlfiber.RateLimit())
	app.Use(fiblogger.New())

	assert := assert.New(t)

	app.Handler()(fctx)

	assert.Equal(http.StatusOK, fctx.Response.StatusCode())
}

func TestRateLimitWithConfigLimited(t *testing.T) {
	t.Parallel()

	cpu, ram, done := rate.MustNewLazy()
	defer close(done)

	time.Sleep(6500 * time.Millisecond)

	limCfg := middlfiber.RateLimit(
		middlfiber.WithConfig(
			middlfiber.Config{
				CommonConfig: middlewares.CommonConfig{
					CPUUtilizationBarrierPercentage: 0.01,
					Logger: golog.New(
						log.New(os.Stderr, "rate/fiber", log.LstdFlags),
					),
				},
			},
		),
		middlfiber.WithLimiter(
			rate.MustNew(
				rate.AppendCPUResourcer(cpu),
				rate.AppendCPUResourcer(ram),
			),
		),
	)

	app, fctx := newApp(limCfg)
	app.Use(fiblogger.New())

	assert := assert.New(t)

	app.Handler()(fctx)

	assert.Equal(http.StatusTooManyRequests, fctx.Response.StatusCode())
	assert.Equal(
		"6",
		string(fctx.Response.Header.Peek(middlewares.HeaderRetryAfter)),
	)
}
