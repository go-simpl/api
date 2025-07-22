package fiberframework

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-simpl/simplapi/pkg/context"
	"github.com/go-simpl/simplapi/pkg/framework"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

type fiberFramework struct {
	app *fiber.App
}

func New(app *fiber.App) framework.Framework {
	return &fiberFramework{
		app: app,
	}
}

func (f *fiberFramework) GET(path string, handler framework.FrameworkHandler) {
	f.app.Get(path, func(c *fiber.Ctx) error {
		return handler(NewRequest(c), NewResponse(c), context.New())
	})
}

func (f *fiberFramework) POST(path string, handler framework.FrameworkHandler) {
	f.app.Post(path, func(c *fiber.Ctx) error {
		return handler(NewRequest(c), NewResponse(c), context.New())
	})
}

func (f *fiberFramework) PUT(path string, handler framework.FrameworkHandler) {
	f.app.Put(path, func(c *fiber.Ctx) error {
		return handler(NewRequest(c), NewResponse(c), context.New())
	})
}

func (f *fiberFramework) PATCH(path string, handler framework.FrameworkHandler) {
	f.app.Patch(path, func(c *fiber.Ctx) error {
		return handler(NewRequest(c), NewResponse(c), context.New())
	})
}

func (f *fiberFramework) DELETE(path string, handler framework.FrameworkHandler) {
	f.app.Delete(path, func(c *fiber.Ctx) error {
		return handler(NewRequest(c), NewResponse(c), context.New())
	})
}

func (f *fiberFramework) OPTIONS(path string, handler framework.FrameworkHandler) {
	f.app.Options(path, func(c *fiber.Ctx) error {
		return handler(NewRequest(c), NewResponse(c), context.New())
	})
}

func (f *fiberFramework) HEAD(path string, handler framework.FrameworkHandler) {
	f.app.Head(path, func(c *fiber.Ctx) error {
		return handler(NewRequest(c), NewResponse(c), context.New())
	})
}

func (f *fiberFramework) TRACE(path string, handler framework.FrameworkHandler) {
	f.app.Trace(path, func(c *fiber.Ctx) error {
		return handler(NewRequest(c), NewResponse(c), context.New())
	})
}

func (f *fiberFramework) ListenAndServe(addr string) error {
	return f.app.Listen(addr)
}

func (f *fiberFramework) TestRequest(req *http.Request) (*http.Response, error) {
	resp := httptest.NewRecorder()
	adaptor.FiberApp(f.app).ServeHTTP(resp, req)
	return resp.Result(), nil
}

func init() {
	framework.RegisterFramework("fiber", func() framework.Framework {
		return New(fiber.New())
	})
}
