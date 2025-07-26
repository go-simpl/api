package fiberframework

import (
	"net/http"
	"net/http/httptest"
	"strings"

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

func (f *fiberFramework) GetNativeApp() interface{} {
	return f.app
}

func (f *fiberFramework) GetOpenAPICompatiblePathPattern(path string) string {
	pathParts := strings.Split(path, "/")
	for i, part := range pathParts {
		if strings.HasPrefix(part, ":") {
			pathParts[i] = "{" + part[1:] + "}"
		}
	}
	return strings.Join(pathParts, "/")
}

func (f *fiberFramework) Register(path string, method string, handler framework.FrameworkHandler) {
	switch method {
	case http.MethodGet:
		f.GET(path, handler)
	case http.MethodPost:
		f.POST(path, handler)
	case http.MethodPut:
		f.PUT(path, handler)
	case http.MethodPatch:
		f.PATCH(path, handler)
	case http.MethodDelete:
		f.DELETE(path, handler)
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

func (f *fiberFramework) ListenAndServe(addr string) error {
	return f.app.Listen(addr)
}

func (f *fiberFramework) Shutdown() error {
	return f.app.Shutdown()
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
