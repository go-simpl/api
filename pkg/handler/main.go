package handler

import (
	"net/http"
	"reflect"

	"github.com/go-simpl/simplapi/errors"
	"github.com/go-simpl/simplapi/pkg/context"
	"github.com/go-simpl/simplapi/pkg/framework"
	"github.com/go-simpl/simplapi/pkg/reflection"
	"github.com/go-simpl/simplapi/types"
)

func WrapHandler(handler interface{}, next framework.FrameworkHandler) framework.FrameworkHandler {
	// First we check if the handler is a function
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		panic("handler must be a function")
	}

	// Parse the function outputs
	numOutputs := handlerType.NumOut()

	if numOutputs == 0 {
		panic("handler must return an error at the least")
	}

	// Last output should be an error
	if handlerType.Out(numOutputs-1).Name() != "error" {
		panic("handler must return an error at the last position")
	}

	createParams := func(req framework.FrameworkRequest, ctx *context.Context) ([]reflect.Value, error) {
		numInputs := handlerType.NumIn()
		inputs := make([]reflect.Value, numInputs)
		for i := 0; i < numInputs; i++ {
			if handlerType.In(i) == reflect.TypeOf(ctx) {
				inputs[i] = reflect.ValueOf(ctx)
				continue
			}

			inputs[i] = reflect.New(handlerType.In(i)).Elem()
			err := reflection.PopulateValueFromTypeUsingContext(req, handlerType.In(i), inputs[i])
			if err != nil {
				return nil, err
			}
		}
		return inputs, nil
	}

	return func(req framework.FrameworkRequest, res framework.FrameworkResponse, ctx *context.Context) error {
		inputs, err := createParams(req, ctx)
		if err != nil {
			if typeErr, ok := err.(errors.TypeError); ok {
				res.SetStatusCode(http.StatusUnprocessableEntity)
				return res.SendJSON(typeErr)
			}
			return err
		}
		results := reflect.ValueOf(handler).Call(inputs)

		// handle error
		errVal := results[len(results)-1]
		if !errVal.IsNil() {
			return errVal.Interface().(error)
		}

		// handle response
		for _, result := range results {
			if !result.IsNil() {
				// check for custom response types
				if response, ok := result.Interface().(*types.HTMLResponse); ok {
					res.SetStatusCode(http.StatusOK)
					res.SetHeader("Content-Type", "text/html")
					return res.SendString(response.HTML)
				}

				if response, ok := result.Interface().(types.APIResponse); ok {
					res.SetStatusCode(response.GetStatusCode())
					return res.SendJSON(result.Interface())
				} else {
					return res.SendJSON(result.Interface())
				}
			}
		}

		if next != nil {
			return next(req, res, ctx)
		}

		return nil
	}
}
