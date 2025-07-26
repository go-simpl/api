package framework

import (
	"mime/multipart"
	"net/http"

	"github.com/go-simpl/simplapi/pkg/context"
)

type Framework interface {
	Register(path string, method string, handler FrameworkHandler)

	ListenAndServe(addr string) error
	Shutdown() error

	GetNativeApp() interface{}

	TestRequest(req *http.Request) (*http.Response, error)
}

type FrameworkHandler func(req FrameworkRequest, res FrameworkResponse, ctx *context.Context) error

type FrameworkRequest interface {
	ParseJSONBody(target interface{}) error
	GetHeader(key string) string
	GetPathParam(key string) string
	GetQueryParam(key string) string
	GetFormValue(key string) string
	GetCookieValue(key string) string
	GetFile(key string) (*multipart.FileHeader, error)
}

type FrameworkResponse interface {
	SetHeader(key string, value string)
	SetCookie(cookie http.Cookie)
	SetStatusCode(statusCode int)
	SendJSON(data interface{}) error
	SendString(data string) error
}
