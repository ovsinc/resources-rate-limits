package fiber

import (
	"log"
	"os"
	"testing"

	"github.com/ovsinc/multilog/golog"
	rate "github.com/ovsinc/resources-rate-limits"
	"github.com/ovsinc/resources-rate-limits/pkg/middlewares"
)

func TestRateLimit(t *testing.T) {
	_ = RateLimit()

	cpu, ram, done := rate.MustNewLazy()
	defer close(done)

	_ = RateLimit(
		WithConfig(
			Config{
				CommonConfig: middlewares.CommonConfig{
					CPUUtilizationBarrierPercentage: 0.01,
					Logger: golog.New(
						log.New(os.Stderr, "rate/fiber ", log.LstdFlags),
					),
				},
			},
		),
		WithLimiter(
			rate.MustNew(
				rate.SetCPUResourcer(cpu),
				rate.SetCPUResourcer(ram),
			),
		),
	)
}
