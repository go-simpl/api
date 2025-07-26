package simplapi

import (
	"github.com/go-simpl/simplapi/types"
)

const swaggerUI = `<!DOCTYPE html>
<html>
<head>
<link type="text/css" rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-themes@3.0.1/themes/3.x/theme-muted.min.css" />
<link rel="shortcut icon" href="https://fastapi.tiangolo.com/img/favicon.png">
<title>Swagger UI</title>
</head>
<body>
<div id="swagger-ui">
</div>
<script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<!-- SwaggerUIBundle is now available on the page -->
<script>
const ui = SwaggerUIBundle({
	url: '/openapi.json',
"dom_id": "#swagger-ui",
"layout": "BaseLayout",
"deepLinking": true,
"showExtensions": true,
"showCommonExtensions": true,
oauth2RedirectUrl: window.location.origin + '/docs/oauth2-redirect',
presets: [
	SwaggerUIBundle.presets.apis,
	SwaggerUIBundle.SwaggerUIStandalonePreset
	],
})
</script>
</body>
</html>
`

const stoplightUI = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>Elements in HTML</title>
  
    <script src="https://unpkg.com/@stoplight/elements/web-components.min.js"></script>
    <link rel="stylesheet" href="https://unpkg.com/@stoplight/elements/styles.min.css">
  </head>
  <body>

    <elements-api
      apiDescriptionUrl="/openapi.json"
      router="hash"
    />

  </body>
</html>`

const redocUI = `<!DOCTYPE html>
<html>
  <head>
    <title>Redoc</title>
    <!-- needed for adaptive design -->
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">

    <!--
    Redoc doesn't change outer page styles
    -->
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <redoc spec-url='/openapi.json'></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"> </script>
  </body>
</html>`

func addOpenAPIRoutes(app *App) {
	app.GET("/_try/swagger", func() (*types.HTMLResponse, error) {
		return &types.HTMLResponse{
			HTML: swaggerUI,
		}, nil
	}).WithoutSpec()

	app.GET("/_try/stoplight", func() (*types.HTMLResponse, error) {
		return &types.HTMLResponse{
			HTML: stoplightUI,
		}, nil
	}).WithoutSpec()

	app.GET("/_try/redoc", func() (*types.HTMLResponse, error) {
		return &types.HTMLResponse{
			HTML: redocUI,
		}, nil
	}).WithoutSpec()

	app.GET("/openapi.json", func() (any, error) {
		return app.spec.ToJson(), nil
	}).WithoutSpec()
}
