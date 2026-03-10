// Package context provides a request-scoped key-value store passed into handlers.
package context

// Context holds data for the duration of a request. Handlers can Set/Get values to pass data along the chain.
type Context struct {
	data map[string]any
}

// New returns a new empty request context.
func New() *Context {
	return &Context{
		data: make(map[string]any),
	}
}

// Set stores a value in the context for this request.
func (c *Context) Set(key string, value any) {
	c.data[key] = value
}

// Get returns the value for key, or nil if not set.
func (c *Context) Get(key string) any {
	return c.data[key]
}

// Contains reports whether key is set in the context.
func (c *Context) Contains(key string) bool {
	_, ok := c.data[key]
	return ok
}

// Delete removes key from the context.
func (c *Context) Delete(key string) {
	delete(c.data, key)
}
