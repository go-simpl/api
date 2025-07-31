package basics

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	simplapi "github.com/go-simpl/simplapi"
	_ "github.com/go-simpl/simplapi/pkg/framework/fiberframework"
	"github.com/go-simpl/simplapi/tests/utils"
)

func TestNoRoutes(t *testing.T) {
	frameworks := utils.GetAllFrameworks()
	for _, framework := range frameworks {
		t.Run(framework, func(t *testing.T) {
			app := simplapi.New(simplapi.WithCreateFramework(framework))
			fApp := app.GetApp()

			req := httptest.NewRequest(http.MethodGet, "/", nil)

			resp, err := fApp.TestRequest(req)
			assert.NoError(t, err)

			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	}
}

type HelloResponse struct {
	Message string `json:"message"`
}

func (r *HelloResponse) GetStatusCode() int {
	return http.StatusOK
}

func TestGET(t *testing.T) {
	frameworks := utils.GetAllFrameworks()
	for _, framework := range frameworks {
		t.Run(framework, func(t *testing.T) {
			app := simplapi.New(simplapi.WithCreateFramework(framework))
			fApp := app.GetApp()

			app.GET("/", func() (*HelloResponse, error) {
				return &HelloResponse{
					Message: "Hello, World!",
				}, nil
			})
			app.Sync()

			req := httptest.NewRequest(http.MethodGet, "/", nil)

			resp, err := fApp.TestRequest(req)
			assert.NoError(t, err)

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")
			respBytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			respString := string(respBytes)
			assert.Equal(t, `{"message":"Hello, World!"}`, respString)
		})
	}
}

func TestListen(t *testing.T) {
	frameworks := utils.GetAllFrameworks()
	for _, framework := range frameworks {
		t.Run(framework, func(t *testing.T) {
			app := simplapi.New(simplapi.WithCreateFramework(framework))

			app.GET("/", func() (*HelloResponse, error) {
				return &HelloResponse{
					Message: "Hello, World!",
				}, nil
			})

			go func() {
				app.ListenAndServe(":3000")
			}()
			defer app.GetApp().Shutdown()

			time.Sleep(1 * time.Second)

			req, err := http.NewRequest(http.MethodGet, "http://localhost:3000/", nil)
			assert.NoError(t, err)
			resp, err := http.DefaultClient.Do(req)

			assert.NoError(t, err)

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")
			respBytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			respString := string(respBytes)
			assert.Equal(t, `{"message":"Hello, World!"}`, respString)
		})
	}
}
