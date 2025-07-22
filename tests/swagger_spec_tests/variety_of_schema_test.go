package swaggerspectests

import (
	"mime/multipart"
	"net/http"
	"testing"

	simplapi "github.com/go-simpl/simplapi"
	_ "github.com/go-simpl/simplapi/pkg/framework/fiberframework"
)

type JsonInput struct {
	ContentType string `header:"Content-Type" example:"application/json"`
	Body        struct {
		Name    string   `json:"name" example:"John Doe"`
		Age     int      `json:"age" example:"25"`
		IsAdmin *bool    `json:"is_admin" example:"false"`
		Rating  float64  `json:"rating" example:"4.5"`
		Tags    []string `json:"tags" example:"[\"tag1\", \"tag2\"]"`
	} `body:"json"`
}

type FormInput struct {
	Body struct {
		Name    string  `form:"name" example:"John Doe"`
		Age     int     `form:"age" example:"25"`
		IsAdmin *bool   `form:"is_admin" example:"false"`
		Rating  float64 `form:"rating" example:"4.5"`
	} `body:"urlencoded"`
}

type MultipartInput struct {
	Body struct {
		Name    string                `form:"name" example:"John Doe"`
		Age     int                   `form:"age" example:"25"`
		IsAdmin *bool                 `form:"is_admin" example:"false"`
		Rating  float64               `form:"rating" example:"4.5"`
		Picture *multipart.FileHeader `form:"picture" example:"picture.jpg"`
	} `body:"multipart"`
}

type Response1 struct {
	Status string `json:"status" example:"OK"`
	Nested []struct {
		Field1 string    `json:"field1" example:"value1"`
		Field2 string    `json:"field2" example:"value2"`
		Random complex64 `json:"random"`
	} `json:"nested"`
}

func (r *Response1) GetStatusCode() int {
	return http.StatusOK
}

type Response2 struct {
	Status string `json:"status" example:"NOT_OK"`
}

func (r *Response2) GetStatusCode() int {
	return http.StatusOK
}

type Response3 struct {
	Status string `json:"status" example:"OK"`
	Nested []struct {
		Field1 string `json:"field1" example:"value1"`
		Field2 string `json:"field2" example:"value2"`
	} `json:"nested"`
}

func (r *Response3) GetStatusCode() int {
	return http.StatusOK
}

func TestAllPossibleSchemaStuff(t *testing.T) {
	app := simplapi.New()

	app.POST("/json", nil, func(input JsonInput) (*Response1, *Response2, *Response3, error) {
		return nil, nil, nil, nil
	})

	app.POST("/form", nil, func(input FormInput) (*Response1, *Response2, *Response3, error) {
		return nil, nil, nil, nil
	})

	app.POST("/multipart", nil, func(input MultipartInput) (*Response1, *Response2, *Response3, error) {
		return nil, nil, nil, nil
	})
}
