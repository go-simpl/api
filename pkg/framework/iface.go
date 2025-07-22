package framework

import (
	"mime/multipart"
	"net/http"
)

type Framework interface {
	GET(path string, handler FrameworkHandler)
	POST(path string, handler FrameworkHandler)
	PUT(path string, handler FrameworkHandler)
	PATCH(path string, handler FrameworkHandler)
	DELETE(path string, handler FrameworkHandler)
	OPTIONS(path string, handler FrameworkHandler)
	HEAD(path string, handler FrameworkHandler)
	TRACE(path string, handler FrameworkHandler)
	ListenAndServe(addr string) error

	TestRequest(req *http.Request) (*http.Response, error)
}

type FrameworkHandler func(req FrameworkRequest, res FrameworkResponse) error

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
	SetStatusCode(statusCode int)
	SendJSON(data interface{}) error
	SendString(data string) error
}
