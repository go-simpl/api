package chaintests

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-simpl/simplapi"

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
