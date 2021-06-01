package fiber_test

import (
	fiber "github.com/gofiber/fiber/v2"
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
