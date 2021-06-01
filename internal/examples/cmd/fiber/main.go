package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ovsinc/multilog/golog"

	ratelimits "github.com/ovsinc/resources-rate-limits"
	"github.com/ovsinc/resources-rate-limits/pkg/middlewares"
	fibermid "github.com/ovsinc/resources-rate-limits/pkg/middlewares/fiber"

	sysfiber "github.com/gofiber/fiber/v2"
	fiblogger "github.com/gofiber/fiber/v2/middleware/logger"
)

var memo []byte = make([]byte, 0)

const size = 20 * 1024 * 1024

type cpuLoad struct {
	count       int
	percent     float64
	timeSeconds int
}

func load(done chan struct{}, lch chan cpuLoad) {
	for {
		select {
		case <-done:
			return
		case l := <-lch:
			RunCPULoad(l.count, l.timeSeconds, int(l.percent))
		}
	}
}

func main() {
	app := sysfiber.New()

	cpu, ram, done := ratelimits.MustNewLazy()
	defer close(done)

	time.Sleep(6 * time.Second)

	app.Use(fiblogger.New())

	app.Use(
		fibermid.RateLimit(
			fibermid.WithConfig(
				fibermid.Config{
					CommonConfig: middlewares.CommonConfig{
						Logger: golog.New(
							log.New(os.Stderr, "rate/fiber", log.LstdFlags),
						),
						Debug: true,
					},
				},
			),
			fibermid.WithLimiter(
				ratelimits.MustNew(
					ratelimits.SetCPUResourcer(cpu),
					ratelimits.SetCPUResourcer(ram),
				),
			),
		),
	)

	app.Get("/", func(c *sysfiber.Ctx) error {
		return c.SendString("Hello tester!")
	})

	app.Get("/ram/add", func(c *sysfiber.Ctx) error {
		memo = append(memo, make([]byte, size)...)
		return c.SendStatus(http.StatusCreated)
	})

	app.Get("/ram/del", func(c *sysfiber.Ctx) error {
		if len(memo) > size {
			memo = memo[:size]
		}
		return c.SendStatus(http.StatusCreated)
	})

	ch := make(chan cpuLoad)
	go load(done, ch)
	defer close(ch)

	app.Get("/cpu/burn/:num/:percent/:time", func(c *sysfiber.Ctx) error {
		num, _ := c.ParamsInt("num")
		percent, _ := c.ParamsInt("percent")
		time, _ := c.ParamsInt("time")

		if num < 1 || percent < 1 || time < 1 {
			return c.SendStatus(http.StatusBadRequest)
		}

		go func() {
			ch <- cpuLoad{
				count:       num,
				percent:     float64(percent),
				timeSeconds: time,
			}
		}()

		return c.SendStatus(http.StatusCreated)
	})

	app.Listen(":8000")
}
