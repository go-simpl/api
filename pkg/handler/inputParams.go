package handler

import (
	"reflect"

	"github.com/go-simpl/simplapi/pkg/context"
	"github.com/go-simpl/simplapi/pkg/framework"
)

func constructParams(req framework.FrameworkRequest, ctx *context.Context, handlerInputTypes []reflect.Type) ([]reflect.Value, error) {
	numInputs := len(handlerInputTypes)
	inputs := make([]reflect.Value, numInputs)
	for i := 0; i < numInputs; i++ {
		if handlerInputTypes[i] == reflect.TypeOf(ctx) {
			inputs[i] = reflect.ValueOf(ctx)
			continue
		}

		inputs[i] = reflect.New(handlerInputTypes[i]).Elem()
		err := computeValuesFromRequest(req, handlerInputTypes[i], inputs[i])
		if err != nil {
			return nil, err
		}
	}
	return inputs, nil
}
