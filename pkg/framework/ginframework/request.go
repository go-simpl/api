package ginframework

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/go-simpl/simplapi/pkg/framework"
)

type ginRequest struct {
	c *gin.Context
}

func NewRequest(c *gin.Context) framework.FrameworkRequest {
	return &ginRequest{
		c: c,
	}
}

func (r *ginRequest) ParseJSONBody(target interface{}) error {
	return r.c.ShouldBindJSON(target)
}

func (r *ginRequest) GetHeader(key string) string {
	return r.c.GetHeader(key)
}

func (r *ginRequest) GetPathParam(key string) string {
	return r.c.Param(key)
}

func (r *ginRequest) GetQueryParam(key string) string {
	return r.c.Query(key)
}

func (r *ginRequest) GetFormValue(key string) string {
	return r.c.PostForm(key)
}

func (r *ginRequest) GetFile(key string) (*multipart.FileHeader, error) {
	return r.c.FormFile(key)
}
