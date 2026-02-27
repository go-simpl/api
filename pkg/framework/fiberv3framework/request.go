package fiberv3framework

import (
	"mime/multipart"

	"github.com/go-simpl/simplapi/pkg/framework"

	"github.com/gofiber/fiber/v3"
)

type fiberV3Request struct {
	c fiber.Ctx
}

func NewRequest(c fiber.Ctx) framework.FrameworkRequest {
	return &fiberV3Request{
		c: c,
	}
}

func (r *fiberV3Request) ParseJSONBody(target interface{}) error {
	return r.c.Bind().Body(target)
}

func (r *fiberV3Request) GetHeader(key string) string {
	return r.c.Get(key)
}

func (r *fiberV3Request) GetPathParam(key string) string {
	return r.c.Params(key)
}

func (r *fiberV3Request) GetQueryParam(key string) string {
	return r.c.Query(key)
}

func (r *fiberV3Request) GetFormValue(key string) string {
	return r.c.FormValue(key)
}

func (r *fiberV3Request) GetCookieValue(key string) string {
	return r.c.Cookies(key)
}

func (r *fiberV3Request) GetFile(key string) (*multipart.FileHeader, error) {
	return r.c.FormFile(key)
}
