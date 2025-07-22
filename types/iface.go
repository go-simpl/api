package types

type APIResponse interface {
	GetStatusCode() int
}

type HTMLResponse struct {
	HTML string
}
