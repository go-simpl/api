package spec

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
)

type Spec struct {
	spec *openapi3.T
	gen  *openapi3gen.Generator
}

func New() *Spec {
	return &Spec{
		spec: &openapi3.T{
			OpenAPI: "3.1.0",
			Info: &openapi3.Info{
				Title:   "SimpleAPI",
				Version: "1.0.0",
			},
			Components: &openapi3.Components{
				Schemas:         openapi3.Schemas{},
				Responses:       openapi3.ResponseBodies{},
				SecuritySchemes: openapi3.SecuritySchemes{},
			},
			Paths: openapi3.NewPaths(),
		},
		gen: openapi3gen.NewGenerator(
			openapi3gen.SchemaCustomizer(func(name string, t reflect.Type, tag reflect.StructTag, schema *openapi3.Schema) error {
				if t.Kind() == reflect.Struct {
					for i := 0; i < t.NumField(); i++ {
						field := t.Field(i)
						if field.Type.Kind() != reflect.Ptr {
							name := field.Tag.Get("json")
							if name != "" {
								schema.Required = append(schema.Required, name)
							}
						}

						if field.Tag.Get("example") != "" {
							name := field.Tag.Get("json")
							if name != "" {
								// TODO: consider typing for the example
								schema.Properties[name].Value.Example = field.Tag.Get("example")
							}
						}
					}
				}
				return nil
			}),
			openapi3gen.CreateComponentSchemas(
				openapi3gen.ExportComponentSchemasOptions{
					ExportComponentSchemas: true,
					ExportTopLevelSchema:   true,
					ExportGenerics:         true,
				},
			),
			openapi3gen.CreateTypeNameGenerator(func(t reflect.Type) string {
				var name string = ""
				if t.Name() == "" {
					name = generateSchemaName()
				} else {
					name = t.Name()
				}

				name = t.PkgPath() + "_" + name

				name = strings.ReplaceAll(name, "github.com/go-simpl/simplapi/example/", "")
				name = strings.ReplaceAll(name, "/", "_")

				return name
			}),
		),
	}
}

func (s *Spec) ToJson() any {
	return s.spec
}

func (s *Spec) Register(path string, method string, tags []string, handlers ...interface{}) {
	if s.spec.Paths.Find(path) == nil {
		s.spec.Paths.Set(path, &openapi3.PathItem{})
	}

	op := openapi3.NewOperation()

	if tags != nil {
		op.Tags = tags
	}

	for _, h := range handlers {
		addRequestInfo(s.gen, s.spec.Components.Schemas, op, h)
	}

	for _, h := range handlers {
		addResponses(s.gen, s.spec.Components.Schemas, op, h)
	}

	switch method {
	case http.MethodGet:
		s.spec.Paths.Value(path).Get = op
	case http.MethodPost:
		s.spec.Paths.Value(path).Post = op
	case http.MethodPut:
		s.spec.Paths.Value(path).Put = op
	case http.MethodDelete:
		s.spec.Paths.Value(path).Delete = op
	case http.MethodPatch:
		s.spec.Paths.Value(path).Patch = op
	}
}
