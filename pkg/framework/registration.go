package framework

// availableFrameworks is populated by each framework's init() (fiber, gin, fiberv3, etc.).
var availableFrameworks map[string]func() Framework

func init() {
	availableFrameworks = map[string]func() Framework{}
}

// RegisterFramework registers a framework by name. Called from framework packages in init().
func RegisterFramework(name string, factory func() Framework) {
	availableFrameworks[name] = factory
}

// GetFramework returns the framework instance for the given name. Panics if not registered.
func GetFramework(name string) Framework {
	factory, ok := availableFrameworks[name]
	if !ok {
		panic("framework not found")
	}
	return factory()
}
