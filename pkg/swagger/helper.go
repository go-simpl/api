package swagger

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-simpl/simplapi/pkg/context"
	"github.com/go-simpl/simplapi/types"
)

func UpdateDefinitionUsingHandler(definition map[string]interface{}, handler interface{}) {
	handlerType := reflect.TypeOf(handler)

	numInputs := handlerType.NumIn()
	paramTypes := make([]reflect.Type, numInputs)

	for i := 0; i < numInputs; i++ {
		paramTypes[i] = handlerType.In(i)
	}

	UpdateDefinitionUsingParamTypes(definition, paramTypes)
}

func UpdateDefinitionUsingParamTypes(definition map[string]interface{}, paramTypes []reflect.Type) {

	HEADER_EXCLUSIONS := map[string]bool{"content-type": true, "content-length": true, "user-agent": true}

	for _, paramType := range paramTypes {
		if paramType == reflect.TypeOf((*context.Context)(nil)) {
			continue
		}

		for i := 0; i < paramType.NumField(); i++ {
			field := paramType.Field(i)
			if field.Tag.Get("body") != "" {
				if field.Tag.Get("body") == "multipart" {
					properties := map[string]interface{}{}

					for j := 0; j < field.Type.NumField(); j++ {
						field := field.Type.Field(j)
						if field.Tag.Get("form") != "" {
							properties[field.Tag.Get("form")] = GetSwaggerSchemaForType(field.Type)
							if field.Tag.Get("example") != "" {
								properties[field.Tag.Get("form")].(map[string]interface{})["example"] = field.Tag.Get("example")
							}
						}
					}

					bodyDefinition := map[string]interface{}{
						"content": map[string]interface{}{
							"multipart/form-data": map[string]interface{}{
								"schema": map[string]interface{}{
									"type":       "object",
									"properties": properties,
								},
							},
						},
					}

					definition["requestBody"] = bodyDefinition

				} else if field.Tag.Get("body") == "urlencoded" {
					properties := map[string]interface{}{}

					for j := 0; j < field.Type.NumField(); j++ {
						field := field.Type.Field(j)
						if field.Tag.Get("form") != "" {
							properties[field.Tag.Get("form")] = GetSwaggerSchemaForType(field.Type)
							if field.Tag.Get("example") != "" {
								properties[field.Tag.Get("form")].(map[string]interface{})["example"] = field.Tag.Get("example")
							}
						}
					}

					bodyDefinition := map[string]interface{}{
						"content": map[string]interface{}{
							"application/x-www-form-urlencoded": map[string]interface{}{
								"schema": map[string]interface{}{
									"type":       "object",
									"properties": properties,
								},
							},
						},
					}

					definition["requestBody"] = bodyDefinition
				} else if field.Tag.Get("body") == "json" {
					bodyDefinition := map[string]interface{}{
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": GetSwaggerSchemaForType(field.Type),
							},
						},
					}
					definition["requestBody"] = bodyDefinition
				}

			} else if field.Type.Kind() == reflect.Struct {
				UpdateDefinitionUsingParamTypes(definition, []reflect.Type{field.Type})

			} else if field.Tag.Get("query") != "" {
				queryDefinition := map[string]interface{}{
					"in":       "query",
					"name":     field.Tag.Get("query"),
					"required": field.Type.Kind() != reflect.Ptr,
					"schema":   GetSwaggerSchemaForType(field.Type),
				}
				definition["parameters"] = append(definition["parameters"].([]interface{}), queryDefinition)

			} else if field.Tag.Get("header") != "" {
				if _, ok := HEADER_EXCLUSIONS[strings.ToLower(field.Tag.Get("header"))]; ok {
					continue
				}
				headerDefinition := map[string]interface{}{
					"in":       "header",
					"name":     field.Tag.Get("header"),
					"required": field.Type.Kind() != reflect.Ptr,
					"schema":   GetSwaggerSchemaForType(field.Type),
				}
				definition["parameters"] = append(definition["parameters"].([]interface{}), headerDefinition)

			} else if field.Tag.Get("path") != "" {
				pathDefinition := map[string]interface{}{
					"in":       "path",
					"name":     field.Tag.Get("path"),
					"required": true, // TODO
					"schema": map[string]interface{}{
						"type": "string",
					},
				}
				definition["parameters"] = append(definition["parameters"].([]interface{}), pathDefinition)
			}
		}
	}
}

func UpdateResponseDefinitionUsingHandler(responses map[string]interface{}, handler interface{}) {
	handlerType := reflect.TypeOf(handler)

	numOutputs := handlerType.NumOut()
	responseTypes := make([]reflect.Type, numOutputs-1)

	for i := 0; i < numOutputs-1; i++ {
		responseTypes[i] = handlerType.Out(i)
	}

	UpdateResponseDefinitionUsingResponseTypes(responses, responseTypes)
}

func UpdateResponseDefinitionUsingResponseTypes(responses map[string]interface{}, responseTypes []reflect.Type) {
	for _, responseType := range responseTypes {
		schema := GetSwaggerSchemaForType(responseType)

		var code int
		var description string

		// Try to create a zero value of the responseType and cast to APIResponse
		if responseType.Kind() == reflect.Ptr {
			val := reflect.New(responseType.Elem()).Interface()
			if apiResp, ok := val.(types.APIResponse); ok {
				code = apiResp.GetStatusCode()
			} else {
				code = 200
			}
		} else {
			panic("response type must be a pointer")
		}

		codeStr := fmt.Sprintf("%d", code)
		description = http.StatusText(code)

		if _, ok := responses[codeStr]; !ok {
			responses[codeStr] = map[string]interface{}{
				"description": description,
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": schema,
					},
				},
			}
		} else {
			// Merge schemas using oneOf
			resp := responses[codeStr].(map[string]interface{})
			resp["description"] = resp["description"].(string) + " or " + description

			content := resp["content"].(map[string]interface{})
			appJson := content["application/json"].(map[string]interface{})
			existingSchema := appJson["schema"]

			// If already oneOf, append; else, create oneOf
			if existingMap, ok := existingSchema.(map[string]interface{}); ok {
				if oneOf, ok := existingMap["oneOf"]; ok {
					existingMap["oneOf"] = append(oneOf.([]interface{}), schema)
				} else {
					appJson["schema"] = map[string]interface{}{
						"oneOf": []interface{}{
							existingSchema,
							schema,
						},
					}
				}
			}
		}
	}
}

func GetSwaggerSchemaForType(t reflect.Type) interface{} {
	fileType := reflect.TypeOf((*multipart.FileHeader)(nil))
	if t.ConvertibleTo(fileType) {
		return map[string]interface{}{
			"type": "file",
		}
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct {
		properties := map[string]interface{}{}

		required := []string{}
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Tag.Get("json") != "" {
				properties[field.Tag.Get("json")] = GetSwaggerSchemaForType(field.Type)
				if field.Tag.Get("example") != "" {
					properties[field.Tag.Get("json")].(map[string]interface{})["example"] = field.Tag.Get("example")
				}

				if field.Type.Kind() != reflect.Ptr {
					required = append(required, field.Tag.Get("json"))
				}
			}
		}

		return map[string]interface{}{
			"type":       "object",
			"properties": properties,
			"required":   required,
		}
	}

	if t.Kind() == reflect.Slice {
		return map[string]interface{}{
			"type":  "array",
			"items": GetSwaggerSchemaForType(t.Elem()),
		}
	}

	switch t.Kind() {
	case reflect.String:
		return map[string]interface{}{
			"type": "string",
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return map[string]interface{}{
			"type": "integer",
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return map[string]interface{}{
			"type": "integer",
		}
	case reflect.Float32, reflect.Float64:
		return map[string]interface{}{
			"type": "number",
		}
	case reflect.Bool:
		return map[string]interface{}{
			"type": "boolean",
		}
	default:
		return map[string]interface{}{
			"type": "string",
		}
	}
}
