package chaintests

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	simplapi "github.com/go-simpl/simplapi"
	"github.com/go-simpl/simplapi/pkg/context"

	_ "github.com/go-simpl/simplapi/pkg/framework/fiberframework"
)

type HelloResponse struct {
	Message string `json:"message"`
}

func (*HelloResponse) GetStatusCode() int {
	return http.StatusOK
}

func TestChainReturnedFromFirstFunc(t *testing.T) {
	app := simplapi.New()
	fApp := app.GetApp()

	func1 := func() (*HelloResponse, error) {
		return &HelloResponse{Message: "func1"}, nil
	}

	func2 := func() (*HelloResponse, error) {
		return &HelloResponse{Message: "func2"}, nil
	}

	app.GET("/", nil, func1, func2)

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	response, err := fApp.TestRequest(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `{"message":"func1"}`, bodyString)
}

func TestChainReturnedFromSecondFunc(t *testing.T) {
	app := simplapi.New()
	fApp := app.GetApp()

	func1 := func() (*HelloResponse, error) {
		return nil, nil
	}

	func2 := func() (*HelloResponse, error) {
		return &HelloResponse{Message: "func2"}, nil
	}

	app.GET("/", nil, func1, func2)

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	response, err := fApp.TestRequest(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `{"message":"func2"}`, bodyString)
}

func TestChainReturnedFromNoFunc(t *testing.T) {
	app := simplapi.New()
	fApp := app.GetApp()

	func1 := func() (*HelloResponse, error) {
		return nil, nil
	}

	func2 := func() (*HelloResponse, error) {
		return nil, nil
	}

	app.GET("/", nil, func1, func2)

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	response, err := fApp.TestRequest(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestChainReturnedErrorFromFirstFunc(t *testing.T) {
	app := simplapi.New()
	fApp := app.GetApp()

	func1 := func() (*HelloResponse, error) {
		return nil, errors.New("server error")
	}

	func2 := func() (*HelloResponse, error) {
		return nil, nil
	}

	app.GET("/", nil, func1, func2)

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	response, err := fApp.TestRequest(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `server error`, bodyString)
}

func TestContextInChain(t *testing.T) {
	app := simplapi.New()
	fApp := app.GetApp()

	func1 := func(ctx *context.Context) (*HelloResponse, error) {
		ctx.Set("test", "test")
		return nil, nil
	}

	func2 := func(ctx *context.Context) (*HelloResponse, error) {
		assert.Equal(t, "test", ctx.Get("test"))
		ctx.Delete("test")
		return nil, nil
	}

	func3 := func(ctx *context.Context) (*HelloResponse, error) {
		assert.False(t, ctx.Contains("test"))
		return &HelloResponse{Message: ""}, nil
	}

	app.GET("/", nil, func1, func2, func3)

	req := httptest.NewRequest("GET", "/", nil)
	response, err := fApp.TestRequest(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
}
