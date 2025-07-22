package injectiontests

import (
	"io"
	"net/http"
	"testing"

	"github.com/go-simpl/simplapi"
	"github.com/stretchr/testify/assert"
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

func doTest[T any](t *testing.T, app *simplapi.App, req *http.Request, routeRegistration RouteRegistrationFunc, pathPattern string) (T, int, string) {
	var result T

	routeRegistration(pathPattern, nil, func(input T) (*HelloResponse, error) {
		result = input
		return &HelloResponse{Message: "Hello"}, nil
	})

	resp, err := app.GetApp().TestRequest(req)
	assert.NoError(t, err)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	return result, resp.StatusCode, string(bodyBytes)
}
