package fiberv3framework

import (
	"net/http"

	"github.com/go-simpl/simplapi/pkg/framework"

	"github.com/gofiber/fiber/v3"
)

type fiberV3Response struct {
	c fiber.Ctx
}

func NewResponse(c fiber.Ctx) framework.FrameworkResponse {
	return &fiberV3Response{
		c: c,
	}
}

func (r *fiberV3Response) SetStatusCode(statusCode int) {
	r.c.Status(statusCode)
}

func (r *fiberV3Response) SetHeader(key string, value string) {
	r.c.Response().Header.Set(key, value)
}

func (r *fiberV3Response) SetCookie(cookie http.Cookie) {
	sameSite := ""
	switch cookie.SameSite {
	case http.SameSiteLaxMode:
		sameSite = "Lax"
	case http.SameSiteStrictMode:
		sameSite = "Strict"
	case http.SameSiteNoneMode:
		sameSite = "None"
	}
	r.c.Cookie(&fiber.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		MaxAge:   cookie.MaxAge,
		Expires:  cookie.Expires,
		Secure:   cookie.Secure,
		HTTPOnly: cookie.HttpOnly,
		SameSite: sameSite,
	})
}

func (r *fiberV3Response) SendJSON(data interface{}) error {
	return r.c.JSON(data)
}

func (r *fiberV3Response) SendString(data string) error {
	return r.c.SendString(data)
}
