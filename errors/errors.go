// Package errors defines error types used by simplapi. TypeError is returned for
// handler input binding failures and is sent as HTTP 422 with a JSON body.
package errors

import (
	"fmt"
)

// TypeError represents a validation or parsing error for a specific field during handler input binding.
type TypeError struct {
	FieldName string `json:"fieldName"`
	Message   string `json:"message"`
}

func (e TypeError) Error() string {
	return fmt.Sprintf("type error: %s: %s", e.FieldName, e.Message)
}
