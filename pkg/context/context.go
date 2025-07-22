package context

type Context struct {
	data map[string]any
}

func New() *Context {
	return &Context{
		data: make(map[string]any),
	}
}

func (c *Context) Set(key string, value any) {
	c.data[key] = value
}

func (c *Context) Get(key string) any {
	return c.data[key]
}

func (c *Context) Contains(key string) bool {
	_, ok := c.data[key]
	return ok
}
