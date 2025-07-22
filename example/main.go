package main

import (
	"github.com/go-simpl/simplapi"

	_ "github.com/go-simpl/simplapi/pkg/framework/fiberframework"
)

func main() {
	app := simplapi.New()
	app.ListenAndServe(":8000")
}
