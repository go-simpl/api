package injectiontests

import (
	"io"
	"net/http"
	"testing"

	simplapi "github.com/go-simpl/simplapi"

	"github.com/stretchr/testify/assert"
)

func SetupTest(frameworkName string, method string) *simplapi.App {
	return simplapi.New(frameworkName, "")
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

func doTest[T any](t *testing.T, app *simplapi.App, req *http.Request, pathPattern string) (T, int, string) {
	var result T

	var registerFunc func(path string, handlers ...interface{}) *simplapi.Endpoint = nil
	switch req.Method {
	case http.MethodGet:
		registerFunc = app.GET
	case http.MethodPost:
		registerFunc = app.POST
	case http.MethodPut:
		registerFunc = app.PUT
	case http.MethodDelete:
		registerFunc = app.DELETE
	case http.MethodPatch:
		registerFunc = app.PATCH
	default:
		panic("unsupported method: " + req.Method)
	}

	registerFunc(pathPattern, func(input T) (*HelloResponse, error) {
		result = input
		return &HelloResponse{Message: "Hello"}, nil
	})
	app.Sync()

	resp, err := app.GetApp().TestRequest(req)
	assert.NoError(t, err)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	return result, resp.StatusCode, string(bodyBytes)
}
