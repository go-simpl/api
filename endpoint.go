package simplapi

// Endpoint is the fluent builder for a single route. Use WithTag, WithSummary, etc.,
// or WithoutSpec to omit it from the OpenAPI spec.
type Endpoint struct {
	method      string
	path        string
	tags        []string
	summary     string
	description string
	operationId string
	handlers    []interface{}
	addToSpec   bool
}

func newEndpoint(method string, path string, handlers ...interface{}) *Endpoint {
	return &Endpoint{
		method:    method,
		path:      path,
		handlers:  handlers,
		addToSpec: true,
	}
}

func (e *Endpoint) WithTag(tag string) *Endpoint {
	e.tags = append(e.tags, tag)
	return e
}

func (e *Endpoint) WithSummary(summary string) *Endpoint {
	e.summary = summary
	return e
}

func (e *Endpoint) WithDescription(description string) *Endpoint {
	e.description = description
	return e
}

func (e *Endpoint) WithOperationId(operationId string) *Endpoint {
	e.operationId = operationId
	return e
}

// WithoutSpec marks this route so it is not added to the OpenAPI spec (e.g. for internal doc routes).
func (e *Endpoint) WithoutSpec() *Endpoint {
	e.addToSpec = false
	return e
}
