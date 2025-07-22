package ginframework

import (
	"github.com/gin-gonic/gin"
	"github.com/go-simpl/simplapi/pkg/framework"
)

type ginResponse struct {
	c *gin.Context
}

func NewResponse(c *gin.Context) framework.FrameworkResponse {
	return &ginResponse{
		c: c,
	}
}

func (r *ginResponse) SetStatusCode(statusCode int) {
	r.c.Status(statusCode)
}

func (r *ginResponse) SetHeader(key string, value string) {
	r.c.Header(key, value)
}

func (r *ginResponse) SendJSON(data interface{}) error {
	r.c.JSON(r.c.Writer.Status(), data)
	return nil
}

func (r *ginResponse) SendString(data string) error {
	r.c.String(r.c.Writer.Status(), data)
	return nil
}
