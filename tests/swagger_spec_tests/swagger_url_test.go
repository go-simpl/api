package swaggerspectests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-simpl/simplapi"
	_ "github.com/go-simpl/simplapi/pkg/framework/fiberframework"
)

func TestOpenAPIJsonURL(t *testing.T) {
	app := simplapi.New()
	fApp := app.GetApp()

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)

	response, err := fApp.TestRequest(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestTryURL(t *testing.T) {
	app := simplapi.New()
	fApp := app.GetApp()

	req := httptest.NewRequest(http.MethodGet, "/try", nil)

	response, err := fApp.TestRequest(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
}
