package fiber_test

import (
	"net/http"
	"testing"
	"time"

	rate "github.com/ovsinc/resources-rate-limits"
	middlfiber "github.com/ovsinc/resources-rate-limits/pkg/middlewares/fiber"
	"github.com/stretchr/testify/require"
)

func BenchmarkFiberWithMiddleware(b *testing.B) {
	cpu, ram, done := rate.MustNewLazy()
	defer close(done)

	time.Sleep(6500 * time.Millisecond)

	limCfg := middlfiber.RateLimit(
		middlfiber.WithLimiter(
			rate.MustNew(
				rate.SetCPUResourcer(cpu),
				rate.SetCPUResourcer(ram),
			),
		),
	)

	app, fctx := newApp(limCfg)

	app.Handler()(fctx)
	require.Equal(b, http.StatusOK, fctx.Response.StatusCode())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.Handler()(fctx)
	}
}

func BenchmarkFiberEmpty(b *testing.B) {
	app, c := newApp(nil)

	app.Handler()(c)
	require.Equal(b, http.StatusOK, c.Response.StatusCode())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.Handler()(c)
	}
}
