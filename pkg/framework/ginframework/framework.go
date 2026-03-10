package ginframework

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-simpl/simplapi/pkg/context"
	"github.com/go-simpl/simplapi/pkg/framework"

	goctx "context"

	"github.com/gin-gonic/gin"
)

type ginFramework struct {
	engine *gin.Engine
	srv    *http.Server
}

func New(engine *gin.Engine) framework.Framework {
	return &ginFramework{
		engine: gin.Default(),
	}
}

func (g *ginFramework) GetNativeApp() interface{} {
	return g.engine
}

// GetOpenAPICompatiblePathPattern returns path unchanged for the OpenAPI spec.
func (g *ginFramework) GetOpenAPICompatiblePathPattern(path string) string {
	return path
}

func (g *ginFramework) Register(path string, method string, handler framework.FrameworkHandler) {
	switch method {
	case http.MethodGet:
		g.GET(path, handler)
	case http.MethodPost:
		g.POST(path, handler)
	case http.MethodPut:
		g.PUT(path, handler)
	case http.MethodPatch:
		g.PATCH(path, handler)
	case http.MethodDelete:
		g.DELETE(path, handler)
	}
}

func (g *ginFramework) GET(path string, handler framework.FrameworkHandler) {
	g.engine.GET(path, func(c *gin.Context) {
		err := handler(NewRequest(c), NewResponse(c), context.New())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	})
}

func (g *ginFramework) POST(path string, handler framework.FrameworkHandler) {
	g.engine.POST(path, func(c *gin.Context) {
		err := handler(NewRequest(c), NewResponse(c), context.New())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	})
}

func (g *ginFramework) PUT(path string, handler framework.FrameworkHandler) {
	g.engine.PUT(path, func(c *gin.Context) {
		err := handler(NewRequest(c), NewResponse(c), context.New())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	})
}

func (g *ginFramework) PATCH(path string, handler framework.FrameworkHandler) {
	g.engine.PATCH(path, func(c *gin.Context) {
		err := handler(NewRequest(c), NewResponse(c), context.New())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	})
}

func (g *ginFramework) DELETE(path string, handler framework.FrameworkHandler) {
	g.engine.DELETE(path, func(c *gin.Context) {
		err := handler(NewRequest(c), NewResponse(c), context.New())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	})
}

func (g *ginFramework) ListenAndServe(addr string) error {
	g.srv = &http.Server{
		Addr:    addr,
		Handler: g.engine,
	}
	return g.srv.ListenAndServe()
}

func (g *ginFramework) Shutdown() error {
	if g.srv == nil {
		return nil
	}
	return g.srv.Shutdown(goctx.Background())
}

func (g *ginFramework) TestRequest(req *http.Request) (*http.Response, error) {
	resp := httptest.NewRecorder()
	g.engine.ServeHTTP(resp, req)
	return resp.Result(), nil
}

func init() {
	framework.RegisterFramework("gin", func() framework.Framework {
		return New(gin.New())
	})
}
