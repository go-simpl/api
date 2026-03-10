// Package simplapi provides a declarative HTTP API builder with optional
// OpenAPI 3.1 spec generation. Use New with options to create an App, then
// register routes via GET, POST, etc. and call ListenAndServe.
package simplapi

import (
	"net/http"

	"github.com/go-simpl/simplapi/pkg/framework"
	"github.com/go-simpl/simplapi/pkg/handler"
	"github.com/go-simpl/simplapi/pkg/spec"
)

// App is the main API builder. It holds a framework adapter, optional OpenAPI
// spec, and a list of endpoints that are registered on Sync or ListenAndServe.
type App struct {
	framework framework.Framework
	spec      *spec.Spec

	endpoints []*Endpoint
}

// AppConfig holds options applied by AppOption functions during New.
type AppConfig struct {
	framework framework.Framework
	spec      *spec.Spec
}

// AppOption configures an App at construction time (e.g. WithCreateFramework, WithAutoOpenAPISpec).
type AppOption func(*AppConfig)

// WithCreateFramework creates an App that uses the named framework (e.g. "fiber", "gin").
func WithCreateFramework(frameworkName string) AppOption {
	return func(c *AppConfig) {
		c.framework = framework.GetFramework(frameworkName)
	}
}

// WithFramework creates an App that uses the given framework instance.
func WithFramework(framework framework.Framework) AppOption {
	return func(c *AppConfig) {
		c.framework = framework
	}
}

// WithAutoOpenAPISpec enables OpenAPI 3.1 spec generation and built-in doc routes. basePkgName is stripped from schema type names.
func WithAutoOpenAPISpec(basePkgName string) AppOption {
	return func(c *AppConfig) {
		c.spec = spec.New(basePkgName)
	}
}

// New creates an App with the given options. If WithAutoOpenAPISpec is used, doc routes are registered immediately.
func New(opts ...AppOption) *App {
	config := &AppConfig{}
	for _, opt := range opts {
		opt(config)
	}

	s := &App{
		framework: config.framework,
		spec:      config.spec,
		endpoints: []*Endpoint{},
	}
	if s.spec != nil {
		addOpenAPIRoutes(s)
	}
	return s
}

// GetApp returns the underlying framework (e.g. for middleware or testing).
func (s *App) GetApp() framework.Framework {
	return s.framework
}

// GET registers a GET route; returns an Endpoint for fluent options (WithTag, WithoutSpec, etc.).
func (s *App) GET(path string, handlers ...interface{}) *Endpoint {
	endpoint := newEndpoint(http.MethodGet, path, handlers...)
	s.addEndpoint(endpoint)
	return endpoint
}

// POST registers a POST route; returns an Endpoint for fluent options.
func (s *App) POST(path string, handlers ...interface{}) *Endpoint {
	endpoint := newEndpoint(http.MethodPost, path, handlers...)
	s.addEndpoint(endpoint)
	return endpoint
}

// PUT registers a PUT route; returns an Endpoint for fluent options.
func (s *App) PUT(path string, handlers ...interface{}) *Endpoint {
	endpoint := newEndpoint(http.MethodPut, path, handlers...)
	s.addEndpoint(endpoint)
	return endpoint
}

// DELETE registers a DELETE route; returns an Endpoint for fluent options.
func (s *App) DELETE(path string, handlers ...interface{}) *Endpoint {
	endpoint := newEndpoint(http.MethodDelete, path, handlers...)
	s.addEndpoint(endpoint)
	return endpoint
}

// PATCH registers a PATCH route; returns an Endpoint for fluent options.
func (s *App) PATCH(path string, handlers ...interface{}) *Endpoint {
	endpoint := newEndpoint(http.MethodPatch, path, handlers...)
	s.addEndpoint(endpoint)
	return endpoint
}

// Sync registers all endpoints with the framework and (if enabled) the OpenAPI spec.
// Registration is deferred until Sync or ListenAndServe so fluent options (WithTag, etc.) are applied first.
func (s *App) Sync() {
	for _, e := range s.endpoints {
		s.framework.Register(e.path, e.method, s.createHandler(e.handlers...))
		if e.addToSpec && s.spec != nil {
			s.spec.Register(s.framework.GetOpenAPICompatiblePathPattern(e.path), e.method, e.tags, e.operationId, e.summary, e.description, e.handlers...)
		}
	}
}

// ListenAndServe registers all endpoints with the framework then starts the server.
func (s *App) ListenAndServe(addr string) error {
	s.Sync()

	return s.framework.ListenAndServe(addr)
}

func (s *App) addEndpoint(endpoint *Endpoint) {
	s.endpoints = append(s.endpoints, endpoint)
}

func (s *App) createHandler(handlers ...interface{}) framework.FrameworkHandler {
	// Reverse so the first handler in the chain is outermost (middleware), last is innermost (route handler).
	for i, j := 0, len(handlers)-1; i < j; i, j = i+1, j-1 {
		handlers[i], handlers[j] = handlers[j], handlers[i]
	}

	var nextHandler framework.FrameworkHandler = nil
	for _, h := range handlers {
		nextHandler = handler.WrapHandler(h, nextHandler)
	}

	return nextHandler
}
