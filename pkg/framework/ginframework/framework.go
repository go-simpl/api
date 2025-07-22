package ginframework

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/go-simpl/simplapi/pkg/context"
	"github.com/go-simpl/simplapi/pkg/framework"
)

type ginFramework struct {
	engine *gin.Engine
}

func New(engine *gin.Engine) framework.Framework {
	return &ginFramework{
		engine: engine,
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

func (g *ginFramework) OPTIONS(path string, handler framework.FrameworkHandler) {
	g.engine.OPTIONS(path, func(c *gin.Context) {
		err := handler(NewRequest(c), NewResponse(c), context.New())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	})
}

func (g *ginFramework) HEAD(path string, handler framework.FrameworkHandler) {
	g.engine.HEAD(path, func(c *gin.Context) {
		err := handler(NewRequest(c), NewResponse(c), context.New())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	})
}

func (g *ginFramework) TRACE(path string, handler framework.FrameworkHandler) {
	g.engine.Handle("TRACE", path, func(c *gin.Context) {
		err := handler(NewRequest(c), NewResponse(c), context.New())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	})
}

func (g *ginFramework) ListenAndServe(addr string) error {
	return g.engine.Run(addr)
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
