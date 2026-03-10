// Package framework defines the abstraction over HTTP frameworks (Fiber, Gin, etc.).
// Implementations register via RegisterFramework in init(); the App uses GetFramework by name.
package framework

import (
	"mime/multipart"
	"net/http"

	"github.com/go-simpl/simplapi/pkg/context"
)

// Framework is the interface that each backend (Fiber, Gin, etc.) must implement.
type Framework interface {
	Register(path string, method string, handler FrameworkHandler)

	ListenAndServe(addr string) error
	Shutdown() error

	GetNativeApp() interface{}
	GetOpenAPICompatiblePathPattern(path string) string

	TestRequest(req *http.Request) (*http.Response, error)
}

// FrameworkHandler is the function type that the framework calls for each request.
type FrameworkHandler func(req FrameworkRequest, res FrameworkResponse, ctx *context.Context) error

// FrameworkRequest abstracts reading body, params, headers, etc. from the underlying framework.
type FrameworkRequest interface {
	ParseJSONBody(target interface{}) error
	GetHeader(key string) string
	GetPathParam(key string) string
	GetQueryParam(key string) string
	GetFormValue(key string) string
	GetCookieValue(key string) string
	GetFile(key string) (*multipart.FileHeader, error)
}

// FrameworkResponse abstracts writing status, headers, and body to the underlying framework.
type FrameworkResponse interface {
	SetHeader(key string, value string)
	SetCookie(cookie http.Cookie)
	SetStatusCode(statusCode int)
	SendJSON(data interface{}) error
	SendString(data string) error
}
