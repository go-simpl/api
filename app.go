package simplapi

import (
	"net/http"

	"github.com/go-simpl/simplapi/pkg/framework"
	"github.com/go-simpl/simplapi/pkg/handler"
	"github.com/go-simpl/simplapi/pkg/spec"
)

type App struct {
	framework framework.Framework
	spec      *spec.Spec

	endpoints []*Endpoint
}

type AppConfig struct {
	framework framework.Framework
	spec      *spec.Spec
}

type AppOption func(*AppConfig)

func WithCreateFramework(frameworkName string) AppOption {
	return func(c *AppConfig) {
		c.framework = framework.GetFramework(frameworkName)
	}
}

func WithFramework(framework framework.Framework) AppOption {
	return func(c *AppConfig) {
		c.framework = framework
	}
}

func WithAutoOpenAPISpec(basePkgName string) AppOption {
	return func(c *AppConfig) {
		c.spec = spec.New(basePkgName)
	}
}

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

func (s *App) GetApp() framework.Framework {
	return s.framework
}

func (s *App) GET(path string, handlers ...interface{}) *Endpoint {
	endpoint := newEndpoint(http.MethodGet, path, handlers...)
	s.addEndpoint(endpoint)
	return endpoint
}

func (s *App) POST(path string, handlers ...interface{}) *Endpoint {
	endpoint := newEndpoint(http.MethodPost, path, handlers...)
	s.addEndpoint(endpoint)
	return endpoint
}

func (s *App) PUT(path string, handlers ...interface{}) *Endpoint {
	endpoint := newEndpoint(http.MethodPut, path, handlers...)
	s.addEndpoint(endpoint)
	return endpoint
}

func (s *App) DELETE(path string, handlers ...interface{}) *Endpoint {
	endpoint := newEndpoint(http.MethodDelete, path, handlers...)
	s.addEndpoint(endpoint)
	return endpoint
}

func (s *App) PATCH(path string, handlers ...interface{}) *Endpoint {
	endpoint := newEndpoint(http.MethodPatch, path, handlers...)
	s.addEndpoint(endpoint)
	return endpoint
}

func (s *App) Sync() {
	for _, e := range s.endpoints {
		s.framework.Register(e.path, e.method, s.createHandler(e.handlers...))
		if e.addToSpec && s.spec != nil {
			s.spec.Register(s.framework.GetOpenAPICompatiblePathPattern(e.path), e.method, e.tags, e.operationId, e.summary, e.description, e.handlers...)
		}
	}
}

func (s *App) ListenAndServe(addr string) error {
	// Actual registration of endpoints happen here
	s.Sync()

	return s.framework.ListenAndServe(addr)
}

func (s *App) addEndpoint(endpoint *Endpoint) {
	s.endpoints = append(s.endpoints, endpoint)
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
