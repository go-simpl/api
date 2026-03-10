package spec

import (
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"github.com/go-simpl/simplapi/pkg/context"
)

// addRequestInfo inspects handler input types and struct tags to add OpenAPI parameters and request body to op.
func addRequestInfo(gen *openapi3gen.Generator, schemas openapi3.Schemas, op *openapi3.Operation, handler interface{}) {
	handlerType := reflect.TypeOf(handler)
	numInputs := handlerType.NumIn()

	for i := 0; i < numInputs; i++ {
		paramType := handlerType.In(i)
		paramName := handlerType.In(i).Name()

		addRequestInfoForParamType(gen, schemas, op, paramType, paramName)
	}
}

// addRequestInfoForParamType adds request info for one param type. context.Context is skipped; body/query/header/cookie/path are mapped to the operation.
func addRequestInfoForParamType(gen *openapi3gen.Generator, schemas openapi3.Schemas, op *openapi3.Operation, paramType reflect.Type, paramName string) {
	if paramType == reflect.TypeOf((*context.Context)(nil)) {
		return
	}

	for i := 0; i < paramType.NumField(); i++ {
		field := paramType.Field(i)
		if field.Tag.Get("body") != "" {
			switch field.Tag.Get("body") {
			case "json":
				schemaRef, _ := gen.NewSchemaRefForValue(reflect.New(field.Type).Interface(), schemas)
				op.RequestBody = &openapi3.RequestBodyRef{
					Value: openapi3.NewRequestBody().WithContent(
						openapi3.NewContentWithJSONSchemaRef(
							schemaRef,
						),
					),
				}
			case "multipart":
				schemaRef := openapi3.NewSchemaRef("", openapi3.NewObjectSchema())
				for j := 0; j < field.Type.NumField(); j++ {
					field := field.Type.Field(j)
					if field.Tag.Get("form") != "" {
						schemaRef.Value.Properties[field.Tag.Get("form")] = openapi3.NewSchemaRef("", getSchemaForType(field.Type))
						if field.Type.Kind() != reflect.Ptr {
							schemaRef.Value.Required = append(schemaRef.Value.Required, field.Tag.Get("form"))
						}
						if field.Tag.Get("example") != "" {
							// TODO: consider typing for the example
							schemaRef.Value.Properties[field.Tag.Get("form")].Value.Example = field.Tag.Get("example")
						}
					}
				}
				op.RequestBody = &openapi3.RequestBodyRef{
					Value: openapi3.NewRequestBody().WithContent(
						openapi3.NewContent(),
					),
				}
				op.RequestBody.Value.Content["multipart/form-data"] = &openapi3.MediaType{
					Schema: schemaRef,
				}
			case "urlencoded":
				schemaRef := openapi3.NewSchemaRef("", openapi3.NewObjectSchema())
				for j := 0; j < field.Type.NumField(); j++ {
					field := field.Type.Field(j)
					if field.Tag.Get("form") != "" {
						schemaRef.Value.Properties[field.Tag.Get("form")] = openapi3.NewSchemaRef("", getSchemaForType(field.Type))
						if field.Type.Kind() != reflect.Ptr {
							schemaRef.Value.Required = append(schemaRef.Value.Required, field.Tag.Get("form"))
						}
						if field.Tag.Get("example") != "" {
							// TODO: consider typing for the example
							schemaRef.Value.Properties[field.Tag.Get("form")].Value.Example = field.Tag.Get("example")
						}
					}
				}
				op.RequestBody = &openapi3.RequestBodyRef{
					Value: openapi3.NewRequestBody().WithContent(
						openapi3.NewContent(),
					),
				}
				op.RequestBody.Value.Content["application/x-www-form-urlencoded"] = &openapi3.MediaType{
					Schema: schemaRef,
				}
			}
		} else if field.Tag.Get("query") != "" {
			op.Parameters = append(
				op.Parameters,
				&openapi3.ParameterRef{
					Value: openapi3.NewQueryParameter(field.Tag.Get("query")).
						WithRequired(field.Type.Kind() != reflect.Ptr).
						WithSchema(getSchemaForType(field.Type)),
				},
			)
		} else if field.Tag.Get("header") != "" {
			op.Parameters = append(
				op.Parameters,
				&openapi3.ParameterRef{
					Value: openapi3.NewHeaderParameter(field.Tag.Get("header")).
						WithRequired(field.Type.Kind() != reflect.Ptr).
						WithSchema(getSchemaForType(field.Type)),
				},
			)
		} else if field.Tag.Get("cookie") != "" {
			op.Parameters = append(
				op.Parameters,
				&openapi3.ParameterRef{
					Value: openapi3.NewCookieParameter(field.Tag.Get("cookie")).
						WithRequired(field.Type.Kind() != reflect.Ptr).
						WithSchema(getSchemaForType(field.Type)),
				},
			)
		} else if field.Tag.Get("path") != "" {
			op.Parameters = append(
				op.Parameters,
				&openapi3.ParameterRef{
					Value: openapi3.NewPathParameter(field.Tag.Get("path")).
						WithRequired(field.Type.Kind() != reflect.Ptr).
						WithSchema(getSchemaForType(field.Type)),
				},
			)
		}
	}
}
