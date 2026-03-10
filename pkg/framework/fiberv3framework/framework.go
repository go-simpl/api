package fiberv3framework

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-simpl/simplapi/pkg/context"
	"github.com/go-simpl/simplapi/pkg/framework"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
)

type fiberV3Framework struct {
	app *fiber.App
}

func New(app *fiber.App) framework.Framework {
	return &fiberV3Framework{
		app: app,
	}
}

func (f *fiberV3Framework) GetNativeApp() interface{} {
	return f.app
}

// GetOpenAPICompatiblePathPattern converts Fiber's :param syntax to OpenAPI's {param} for the spec.
func (f *fiberV3Framework) GetOpenAPICompatiblePathPattern(path string) string {
	pathParts := strings.Split(path, "/")
	for i, part := range pathParts {
		if strings.HasPrefix(part, ":") {
			pathParts[i] = "{" + part[1:] + "}"
		}
	}
	return strings.Join(pathParts, "/")
}

func (f *fiberV3Framework) Register(path string, method string, handler framework.FrameworkHandler) {
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

func (f *fiberV3Framework) GET(path string, handler framework.FrameworkHandler) {
	f.app.Get(path, func(c fiber.Ctx) error {
		return handler(NewRequest(c), NewResponse(c), context.New())
	})
}

func (f *fiberV3Framework) POST(path string, handler framework.FrameworkHandler) {
	f.app.Post(path, func(c fiber.Ctx) error {
		return handler(NewRequest(c), NewResponse(c), context.New())
	})
}

func (f *fiberV3Framework) PUT(path string, handler framework.FrameworkHandler) {
	f.app.Put(path, func(c fiber.Ctx) error {
		return handler(NewRequest(c), NewResponse(c), context.New())
	})
}

func (f *fiberV3Framework) PATCH(path string, handler framework.FrameworkHandler) {
	f.app.Patch(path, func(c fiber.Ctx) error {
		return handler(NewRequest(c), NewResponse(c), context.New())
	})
}

func (f *fiberV3Framework) DELETE(path string, handler framework.FrameworkHandler) {
	f.app.Delete(path, func(c fiber.Ctx) error {
		return handler(NewRequest(c), NewResponse(c), context.New())
	})
}

func (f *fiberV3Framework) ListenAndServe(addr string) error {
	return f.app.Listen(addr)
}

func (f *fiberV3Framework) Shutdown() error {
	return f.app.Shutdown()
}

func (f *fiberV3Framework) TestRequest(req *http.Request) (*http.Response, error) {
	resp := httptest.NewRecorder()
	adaptor.FiberApp(f.app).ServeHTTP(resp, req)
	return resp.Result(), nil
}

func init() {
	framework.RegisterFramework("fiberv3", func() framework.Framework {
		return New(fiber.New())
	})
}
