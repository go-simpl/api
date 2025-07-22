package fiberframework

import (
	"mime/multipart"

	"github.com/go-simpl/simplapi/pkg/framework"
	"github.com/gofiber/fiber/v2"
)

type fiberRequest struct {
	c *fiber.Ctx
}

func NewRequest(c *fiber.Ctx) framework.FrameworkRequest {
	return &fiberRequest{
		c: c,
	}
}

func (r *fiberRequest) ParseJSONBody(target interface{}) error {
	return r.c.BodyParser(target)
}

func (r *fiberRequest) GetHeader(key string) string {
	return r.c.Get(key)
}

func (r *fiberRequest) GetPathParam(key string) string {
	return r.c.Params(key)
}

func (r *fiberRequest) GetQueryParam(key string) string {
	return r.c.Query(key)
}

func (r *fiberRequest) GetFormValue(key string) string {
	return r.c.FormValue(key)
}

func (r *fiberRequest) GetCookieValue(key string) string {
	return r.c.Cookies(key)
}

func (r *fiberRequest) GetFile(key string) (*multipart.FileHeader, error) {
	return r.c.FormFile(key)
}
