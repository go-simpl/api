package injectiontests

import (
	"net/http"

	"github.com/go-simpl/simplapi"
)

type RouteRegistrationFunc func(path string, tags []string, handlers ...interface{})

func SetupTest(frameworkName string, method string) (*simplapi.App, RouteRegistrationFunc) {
	app := simplapi.New(frameworkName)
	switch method {
	case http.MethodGet:
		return app, app.GET
	case http.MethodPost:
		return app, app.POST
	case http.MethodPut:
		return app, app.PUT
	case http.MethodDelete:
		return app, app.DELETE
	case http.MethodPatch:
		return app, app.PATCH
	default:
		panic("Invalid method")
	}
}

type HelloResponse struct {
	Message string `json:"message"`
}

func shouldSkipOptionalPathParamTest(frameworkName string) bool {
	switch frameworkName {
	case "gin":
		return true
	}
	return false
}
