package handler

import (
	"net/http"
	"reflect"

	"github.com/go-simpl/simplapi/errors"
	"github.com/go-simpl/simplapi/pkg/context"
	"github.com/go-simpl/simplapi/pkg/framework"
)

// WrapHandler wraps a user handler (function) into a FrameworkHandler. The handler's last return must be error;
// the first non-nil return value is sent as the response. If no response is sent, next is called.
func WrapHandler(handler interface{}, next framework.FrameworkHandler) framework.FrameworkHandler {
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		panic("handler must be a function")
	}

	numOutputs := handlerType.NumOut()

	if numOutputs == 0 {
		panic("handler must return an error at the least")
	}

	if handlerType.Out(numOutputs-1).Name() != "error" {
		panic("handler must return an error at the last position")
	}

	handlerInputTypes := make([]reflect.Type, handlerType.NumIn())
	for i := 0; i < handlerType.NumIn(); i++ {
		handlerInputTypes[i] = handlerType.In(i)
	}

	return func(req framework.FrameworkRequest, res framework.FrameworkResponse, ctx *context.Context) error {
		inputs, err := constructParams(req, ctx, handlerInputTypes)
		if err != nil {
			if typeErr, ok := err.(errors.TypeError); ok {
				res.SetStatusCode(http.StatusUnprocessableEntity)
				return res.SendJSON(typeErr)
			}
			return err
		}
		results := reflect.ValueOf(handler).Call(inputs)

		errVal := results[len(results)-1]
		if !errVal.IsNil() {
			return errVal.Interface().(error)
		}

		for _, result := range results {
			if !result.IsNil() {
				bodyDone, err := transferToResponse(res, result)
				if err != nil {
					return err
				}

				if bodyDone {
					return nil
				}
			}
		}

		if next != nil {
			return next(req, res, ctx)
		}

		return nil
	}
}
