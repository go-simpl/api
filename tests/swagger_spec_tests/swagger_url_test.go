package swaggerspectests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	simplapi "github.com/go-simpl/simplapi"
	_ "github.com/go-simpl/simplapi/pkg/framework/fiberframework"
	"github.com/go-simpl/simplapi/tests/utils"
)

func TestOpenAPIJsonURL(t *testing.T) {
	frameworks := utils.GetAllFrameworks()
	for _, framework := range frameworks {
		app := simplapi.New(simplapi.WithCreateFramework(framework), simplapi.WithAutoOpenAPISpec("tests"))
		app.Sync()
		fApp := app.GetApp()

		req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)

		response, err := fApp.TestRequest(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
}

func TestTryURL(t *testing.T) {
	frameworks := utils.GetAllFrameworks()
	for _, framework := range frameworks {
		app := simplapi.New(simplapi.WithCreateFramework(framework), simplapi.WithAutoOpenAPISpec("tests"))
		app.Sync()
		fApp := app.GetApp()
		{
			req := httptest.NewRequest(http.MethodGet, "/_try/stoplight", nil)

			response, err := fApp.TestRequest(req)
			assert.NoError(t, err)

			assert.Equal(t, http.StatusOK, response.StatusCode)
		}
		{
			req := httptest.NewRequest(http.MethodGet, "/_try/swagger", nil)

			response, err := fApp.TestRequest(req)
			assert.NoError(t, err)

			assert.Equal(t, http.StatusOK, response.StatusCode)
		}
		{
			req := httptest.NewRequest(http.MethodGet, "/_try/redoc", nil)

			response, err := fApp.TestRequest(req)
			assert.NoError(t, err)

			assert.Equal(t, http.StatusOK, response.StatusCode)
		}
	}
}
