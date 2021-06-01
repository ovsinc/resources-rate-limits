// +build os

package fiber_test

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	fiblogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ovsinc/multilog/golog"
	rate "github.com/ovsinc/resources-rate-limits"
	"github.com/ovsinc/resources-rate-limits/pkg/middlewares"
	middlfiber "github.com/ovsinc/resources-rate-limits/pkg/middlewares/fiber"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitDefaultOK(t *testing.T) {
	t.Parallel()

	_, ram := rate.MustNewSimple()

	app, fctx := newApp(
		middlfiber.RateLimit(
			middlfiber.WithConfig(
				middlfiber.Config{
					CommonConfig: middlewares.CommonConfig{
						Logger: golog.New(
							log.New(os.Stderr, "rate/fiber ", log.LstdFlags),
						),
						Debug: true,
					},
				},
			),
			middlfiber.WithLimiter(
				rate.MustNew(
					rate.SetCPUResourcer(ram),
				),
			),
		),
	)

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
						log.New(os.Stderr, "rate/fiber ", log.LstdFlags),
					),
				},
			},
		),
		middlfiber.WithLimiter(
			rate.MustNew(
				rate.SetCPUResourcer(cpu),
				rate.SetCPUResourcer(ram),
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

func TestRateLimitWithConfigClosed(t *testing.T) {
	t.Parallel()

	cpu, ram, done := rate.MustNewLazy()
	close(done)

	time.Sleep(300 * time.Millisecond)

	limCfg := middlfiber.RateLimit(
		middlfiber.WithConfig(
			middlfiber.Config{
				CommonConfig: middlewares.CommonConfig{
					CPUUtilizationBarrierPercentage: 0.01,
					Logger: golog.New(
						log.New(os.Stderr, "rate/fiber ", log.LstdFlags),
					),
				},
			},
		),
		middlfiber.WithLimiter(
			rate.MustNew(
				rate.SetCPUResourcer(cpu),
				rate.SetCPUResourcer(ram),
			),
		),
	)

	app, fctx := newApp(limCfg)
	app.Use(fiblogger.New())

	assert := assert.New(t)

	app.Handler()(fctx)

	assert.Equal(http.StatusTooManyRequests, fctx.Response.StatusCode())
	assert.Equal(
		"3",
		string(fctx.Response.Header.Peek(middlewares.HeaderRetryAfter)),
	)
}
