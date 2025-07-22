package framework

var availableFrameworks map[string]func() Framework

func init() {
	availableFrameworks = map[string]func() Framework{}
}

func RegisterFramework(name string, factory func() Framework) {
	availableFrameworks[name] = factory
}

func GetFramework(name string) Framework {
	factory, ok := availableFrameworks[name]
	if !ok {
		panic("framework not found")
	}
	return factory()
}
