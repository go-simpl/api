package errors

import (
	"fmt"
)

type TypeError struct {
	FieldName string `json:"fieldName"`
	Message   string `json:"message"`
}

func (e TypeError) Error() string {
	return fmt.Sprintf("type error: %s: %s", e.FieldName, e.Message)
}
