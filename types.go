package simplapi

// Q wraps a value of any type T (e.g. for generic handler or spec use).
type Q[T any] struct {
	Value T
}
