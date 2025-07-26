# go-simpl/api - Simple yet powerful API framework for GO (with swagger spec generation)
![Coverage](https://img.shields.io/badge/Coverage-88.9%25-brightgreen)

A simple, developer-friendly API framework for Go, featuring automatic Swagger (OpenAPI) spec generation.

## 🚀 Features

- **Intuitive handler definitions** using Go structs
- **Automatic request parsing** (JSON, form, multipart, etc.)
- **Zero-boilerplate route registration** for all HTTP methods
- **Instant Swagger docs** generated from your handlers

## 🛠️ Quick Start

**Create a new app**

```go
import (
    "github.com/go-simpl/simplapi"

    _ "github.com/go-simpl/simplapi/pkg/framework/fiberframework"
)

func main() {
    app := simplapi.New("fiber")

    // Add routes

    app.ListenAndServe(":8000")
}
```

Swagger UI will be instantly available on `http://localhost:8000/try` and the swagger (OpenAPI) spec on `http://localhost:8000/openapi.json`

**Adding routes**

To add routes, call the appropriate methods on the app.

```go
app.GET(path, handlers...)
app.POST(path, handlers...)
app.PUT(path, handlers...)
// and so on
```

path: `string` - path to handle

handlers: `[]interface{}` - functions that handle the request

### Handler

Handler function is the important building block that controls how the endpoints are processed and the specs are generated. The function should be of the following format:

```go
func HandlerFunc(input InputType) (*OutputType1, *OutputType2, error) {
    // Handler code
}
```

The function must return one of the output types or error. The function must, at the least, return error.

Output type will be returned as `JSON` response by marshalling the struct.

Additionally, if the Output type defines GetStatusCode function, that will be used to set the status of the response.

If you provide multiple handler functions for a route, they are called one after another in the order you specify. Each handler is executed until one of them returns a non-nil output or an error. As soon as a handler returns a result (any non-nil output or error), the chain stops and no further handlers are called. If a handler returns only nil values, the next handler in the list will be executed.

#### InputType definition

Input for the API can be defined by struct field tags.

**Query Params**

```go
type InputType struct {
    PageNum int `query:"page_num"`
    PageSize int `query:"page_size"`
}
```

**Path Params**

```go
// Assuming path is `/books/:book_id`
type InputType struct {
    BookId string `path:"book_id"`
}
```

**Headers**

```go
type InputType struct {
    ApiKey string `header:"X-API-Key"`
}
```

**JSON Body**

```go
type InputType struct {
    Body struct {
        Title string `json:"title"`
        // ... JSON fields here
    } `body:"json"`
}
```

**URL Encoded from**

```go
type InputType struct {
    Form struct {
        Email    string `form:"email"`
        Password string `form:"password"`
    } `body:"urlencoded"`
}
```

**Multipart form**

```go
type InputType struct {
    Form struct {
        UploadedFile *multipart.FileHeader `form:"file"`
    } `body:"multipart"`
}
```
