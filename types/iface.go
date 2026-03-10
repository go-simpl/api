// Package types defines handler return types that drive response status and body.
package types

// APIResponse is implemented by structs returned from handlers to set status code and JSON body.
type APIResponse interface {
	GetStatusCode() int
}

// HTMLResponse is returned from handlers to send HTML with Content-Type text/html.
type HTMLResponse struct {
	HTML string
}
