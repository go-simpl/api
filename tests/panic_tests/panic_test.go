package panictests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	simplapi "github.com/go-simpl/simplapi"
	_ "github.com/go-simpl/simplapi/pkg/framework/fiberframework"
)

type HelloResponse struct {
	Message string `json:"message"`
}

func (*HelloResponse) GetStatusCode() int {
	return http.StatusOK
}

func TestNonFuncHandler(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		assert.Equal(t, "handler must be a function", r)
	}()
	app := simplapi.New()
	app.GET("/", nil, "hello")
}

func TestNoOutputHandler(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		assert.Equal(t, "handler must return an error at the least", r)
	}()
	app := simplapi.New()
	app.GET("/", nil, func() {})
}

func TestHandlerThatDoesnReturnError(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		assert.Equal(t, "handler must return an error at the last position", r)
	}()
	app := simplapi.New()
	app.GET("/", nil, func() *HelloResponse {
		return nil
	})
}

func TestHandlerResponseMustBePointer(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		assert.Equal(t, "response type must be a pointer", r)
	}()
	app := simplapi.New()
	app.GET("/", nil, func() (HelloResponse, error) {
		return HelloResponse{}, nil
	})
}
