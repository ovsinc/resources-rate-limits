package fiber_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	rate "gitlab.com/ovsinc/memory-rate-limits"
	middlfiber "gitlab.com/ovsinc/memory-rate-limits/middlewares/fiber"
	memmock "gitlab.com/ovsinc/memory-rate-limits/resources/memory/mock"
)

func BenchmarkFiberWithMiddleware(b *testing.B) {
	cfg := middlfiber.DefaultFiberConfig
	cfg.Rate = rate.New(
		rate.WithMemoryChecker(memmock.NewMemUnlimited()),
	)

	app, c := newApp(middlfiber.RateLimitWithConfig(&cfg))

	app.Handler()(c)
	require.Equal(b, http.StatusOK, c.Response.StatusCode())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.Handler()(c)
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
