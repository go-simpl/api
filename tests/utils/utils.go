package utils

import (
	_ "github.com/go-simpl/simplapi/pkg/framework/fiberframework"
	_ "github.com/go-simpl/simplapi/pkg/framework/fiberv3framework"
	_ "github.com/go-simpl/simplapi/pkg/framework/ginframework"
)

func GetAllFrameworks() []string {
	return []string{
		"fiber",
		"fiberv3",
		"gin",
	}
}

func GetSupportedTypes() []string {
	return []string{
		"string",
		"int",
		"uint",
		"float64",
		"bool",
	}
}
