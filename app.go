package simplapi

import (
	"github.com/go-simpl/simplapi/pkg/framework"
	"github.com/go-simpl/simplapi/pkg/handler"
	"github.com/go-simpl/simplapi/pkg/swagger"
)

type App struct {
	framework   framework.Framework
	swaggerJson map[string]interface{}
}

func New(frameworkName ...string) *App {
	frameworkName = append(frameworkName, "fiber")

	s := &App{
		framework: framework.GetFramework(frameworkName[0]),
		swaggerJson: map[string]interface{}{
			"openapi": "3.0.0",
			"info": map[string]interface{}{
				"title":   "SimpleAPI",
				"version": "1.0.0",
			},
			"paths": map[string]interface{}{},
		},
	}
	addSwaggerRoutes(s)
	return s
}

func (s *App) GetApp() framework.Framework {
	return s.framework
}

func (s *App) ListenAndServe(addr string) error {
	return s.framework.ListenAndServe(addr)
}

func (s *App) createHandler(handlers ...interface{}) framework.FrameworkHandler {
	// Reverse the handlers slice
	for i, j := 0, len(handlers)-1; i < j; i, j = i+1, j-1 {
		handlers[i], handlers[j] = handlers[j], handlers[i]
	}

	var nextHandler framework.FrameworkHandler = nil
	for _, h := range handlers {
		nextHandler = handler.WrapHandler(h, nextHandler)
	}

	return nextHandler
}

func (s *App) GET(path string, tags []string, handlers ...interface{}) {
	s.framework.GET(path, s.createHandler(handlers...))
	s.addToSwagger(path, "get", handlers, tags)
}

func (s *App) POST(path string, tags []string, handlers ...interface{}) {
	s.framework.POST(path, s.createHandler(handlers...))
	s.addToSwagger(path, "post", handlers, tags)
}

func (s *App) PUT(path string, tags []string, handlers ...interface{}) {
	s.framework.PUT(path, s.createHandler(handlers...))
	s.addToSwagger(path, "put", handlers, tags)
}

func (s *App) PATCH(path string, tags []string, handlers ...interface{}) {
	s.framework.PATCH(path, s.createHandler(handlers...))
	s.addToSwagger(path, "patch", handlers, tags)
}

func (s *App) DELETE(path string, tags []string, handlers ...interface{}) {
	s.framework.DELETE(path, s.createHandler(handlers...))
	s.addToSwagger(path, "delete", handlers, tags)
}

func (s *App) addToSwagger(path string, method string, handlers []interface{}, tags []string) {
	if path == "/try" || path == "/openapi.json" {
		return
	}

	definition := map[string]interface{}{
		"parameters": []interface{}{},
		"responses":  map[string]interface{}{},
		"tags":       tags,
	}

	if method != "get" {
		definition["requestBody"] = map[string]interface{}{}
	}

	for _, handler := range handlers {
		swagger.UpdateDefinitionUsingHandler(definition, handler)
	}

	if _, ok := s.swaggerJson["paths"].(map[string]interface{})[path]; !ok {
		s.swaggerJson["paths"].(map[string]interface{})[path] = map[string]interface{}{}
	}

	if _, ok := s.swaggerJson["paths"].(map[string]interface{})[path].(map[string]interface{})[method]; !ok {
		s.swaggerJson["paths"].(map[string]interface{})[path].(map[string]interface{})[method] = definition
	}

	responses := definition["responses"].(map[string]interface{})
	for _, handler := range handlers {
		swagger.UpdateResponseDefinitionUsingHandler(responses, handler)
	}
}
