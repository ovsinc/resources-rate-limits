package main

import (
	"net/http"

	fibermid "github.com/ovsinc/resources-rate-limits/pkg/middlewares/fiber"

	sysfiber "github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
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

	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)

	app.Use(fibermid.RateLimitWithConfig(&fibermid.DefaultFiberConfig))

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

	done := make(chan struct{})
	ch := make(chan cpuLoad)
	go load(done, ch)
	defer func() {
		close(done)
		close(ch)
	}()

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
