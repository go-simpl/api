package utils

import (
	_ "github.com/go-simpl/simplapi/pkg/framework/fiberframework"
	_ "github.com/go-simpl/simplapi/pkg/framework/ginframework"
)

func GetAllFrameworks() []string {
	return []string{
		"fiber",
		"gin",
	}
}
