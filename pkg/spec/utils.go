package spec

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
)

func getSchemaForType(t reflect.Type) *openapi3.Schema {
	if t == reflect.TypeOf((*multipart.FileHeader)(nil)) {
		schema := openapi3.NewSchema()
		schema.Type = &openapi3.Types{"file"}
		return schema
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return openapi3.NewIntegerSchema()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return openapi3.NewIntegerSchema()
	case reflect.Float32, reflect.Float64:
		return openapi3.NewFloat64Schema()
	case reflect.Bool:
		return openapi3.NewBoolSchema()
	case reflect.String:
		return openapi3.NewStringSchema()
	}

	return openapi3.NewStringSchema()
}

var schemaNameCount int = 0
var schemaNameMutex sync.Mutex

func generateSchemaName() string {
	schemaNameMutex.Lock()
	defer schemaNameMutex.Unlock()
	schemaNameCount++
	return fmt.Sprintf("Schema%d", schemaNameCount)
}
