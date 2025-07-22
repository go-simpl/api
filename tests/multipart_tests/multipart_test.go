package multiparttests

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-simpl/simplapi"
	_ "github.com/go-simpl/simplapi/pkg/framework/fiberframework"
)

func TestFileUpload(t *testing.T) {
	app := simplapi.New()
	fApp := app.GetApp()

	type FileUploadInput struct {
		Body struct {
			File *multipart.FileHeader `form:"file"`
		} `body:"multipart"`
	}

	app.POST("/upload", nil, func(input FileUploadInput) error {
		assert.Equal(t, "test.txt", input.Body.File.Filename)
		f, err := input.Body.File.Open()
		assert.NoError(t, err)
		fBytes, err := io.ReadAll(f)
		assert.NoError(t, err)
		assert.Equal(t, "Hello, world!", string(fBytes))
		return nil
	})

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, err := writer.CreateFormFile("file", "test.txt")
	assert.NoError(t, err)
	_, err = file.Write([]byte("Hello, world!"))
	assert.NoError(t, err)
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	response, err := fApp.TestRequest(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
}
