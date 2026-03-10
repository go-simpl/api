package handler

import (
	"net/http"
	"reflect"

	"github.com/go-simpl/simplapi/pkg/framework"
	"github.com/go-simpl/simplapi/types"
)

// hasResponseTags is reserved for future response tagging; currently always false.
func hasResponseTags(t reflect.Type) bool {
	return false
}

func transferToResponse(res framework.FrameworkResponse, result reflect.Value) (bool, error) {
	if hasResponseTags(result.Type()) {
		return false, nil
	}

	// First matching type wins: HTMLResponse (text/html), then APIResponse (status + JSON), else raw JSON.
	if response, ok := result.Interface().(*types.HTMLResponse); ok {
		res.SetStatusCode(http.StatusOK)
		res.SetHeader("Content-Type", "text/html")
		return true, res.SendString(response.HTML)
	}

	if response, ok := result.Interface().(types.APIResponse); ok {
		res.SetStatusCode(response.GetStatusCode())
		return true, res.SendJSON(result.Interface())
	} else {
		return true, res.SendJSON(result.Interface())
	}
}
