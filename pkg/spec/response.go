package spec

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"github.com/go-simpl/simplapi/types"
)

func addResponses(gen *openapi3gen.Generator, schemas openapi3.Schemas, op *openapi3.Operation, handler interface{}) {
	handlerType := reflect.TypeOf(handler)

	numOutputs := handlerType.NumOut()

	for i := 0; i < numOutputs-1; i++ {
		addResponseTypeToOperation(gen, schemas, op, handlerType.Out(i))
	}
}

func addResponseTypeToOperation(gen *openapi3gen.Generator, schemas openapi3.Schemas, op *openapi3.Operation, responseType reflect.Type) {
	if op.Responses == nil {
		op.Responses = openapi3.NewResponses()
	}

	if responseType.Kind() != reflect.Ptr {
		panic("response type must be a pointer")
	}

	var code int = 200

	if apiResp, ok := reflect.New(responseType.Elem()).Interface().(types.APIResponse); ok {
		code = apiResp.GetStatusCode()
	}

	schemaRef, _ := gen.NewSchemaRefForValue(reflect.New(responseType.Elem()).Interface(), schemas)

	codeStr := fmt.Sprintf("%d", code)

	if op.Responses.Value(codeStr) == nil {
		response := openapi3.NewResponse().WithDescription(http.StatusText(code))
		response.Content = openapi3.NewContentWithJSONSchemaRef(schemaRef)

		op.Responses.Set(fmt.Sprintf("%d", code), &openapi3.ResponseRef{
			Value: response,
		})
		return
	}

	// Response already exists, add to oneOf if oneOf exists or create it
	response := op.Responses.Value(codeStr).Value
	content := response.Content
	existingSchemaRef := content.Get("application/json").Schema
	if existingSchemaRef.Value == nil {
		// One of hasn't been created here, so create and move current ref into it
		existingSchemaRef.Value = openapi3.NewSchema()
		existingSchemaRef.Value.OneOf = append(existingSchemaRef.Value.OneOf, openapi3.NewSchemaRef(existingSchemaRef.Ref, nil))
		existingSchemaRef.Ref = ""
	}
	existingSchemaRef.Value.OneOf = append(existingSchemaRef.Value.OneOf, schemaRef)
}
