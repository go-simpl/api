package injectiontests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-simpl/simplapi"
	_ "github.com/go-simpl/simplapi/pkg/framework/fiberframework"
)

func TestPathParam(t *testing.T) {
	frameworks := []string{
		"fiber",
	}
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		// http.MethodHead,
		// http.MethodOptions,
		// http.MethodTrace,
	}

	for _, framework := range frameworks {
		t.Run(framework, func(t *testing.T) {
			for _, method := range methods {
				t.Run(method, func(t *testing.T) {
					t.Run("string", func(t *testing.T) { testPathParamOfTypeString(t, method) })
					t.Run("int", func(t *testing.T) { testPathParamOfTypeInt(t, method) })
					t.Run("uint", func(t *testing.T) { testPathParamOfTypeUint(t, method) })
					t.Run("float", func(t *testing.T) { testPathParamOfTypeFloat(t, method) })
					t.Run("bool", func(t *testing.T) { testPathParamOfTypeBool(t, method) })
				})
			}
		})
	}
}

func doPathTest[T any](t *testing.T, app *simplapi.App, req *http.Request, routeRegistration RouteRegistrationFunc, pathPattern string) (T, int, string) {
	var result T

	routeRegistration(pathPattern, nil, func(input T) (*HelloResponse, error) {
		result = input
		return &HelloResponse{Message: "Hello"}, nil
	})

	resp, err := app.GetApp().TestRequest(req)
	assert.NoError(t, err)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	return result, resp.StatusCode, string(bodyBytes)
}

func testPathParamOfTypeString(t *testing.T, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Name string `path:"name"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/hello/John", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/hello/:name")
				assert.Equal(t, "John", res.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Name string `path:"name"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/hello/", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/hello/:name")
				assert.Equal(t, "", res.Name)
				assert.Equal(t, http.StatusNotFound, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name string `path:"name"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/hello/John", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/hello/:name")
				assert.Equal(t, "John", res.Inner.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name string `path:"name"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/hello/", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/hello/:name")
				assert.Equal(t, "", res.Inner.Name)
				assert.Equal(t, http.StatusNotFound, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Name *string `path:"name"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/hello/John", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/hello/:name?")
				assert.NotNil(t, res.Name)
				assert.Equal(t, "John", *res.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Name *string `path:"name"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/hello/", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/hello/:name?")
				assert.Nil(t, res.Name)
				assert.Equal(t, http.StatusOK, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name *string `path:"name"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/hello/John", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/hello/:name?")
				assert.NotNil(t, res.Inner.Name)
				assert.Equal(t, "John", *res.Inner.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name *string `path:"name"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/hello/", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/hello/:name?")
				assert.Nil(t, res.Inner.Name)
				assert.Equal(t, http.StatusOK, status)
			})
		})
	})
}

func testPathParamOfTypeInt(t *testing.T, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Age int `path:"age"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/user/25", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/user/:age")
				assert.Equal(t, 25, res.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Age int `path:"age"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/user/invalid", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/user/:age")
				assert.Equal(t, 0, res.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age int `path:"age"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/user/25", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/user/:age")
				assert.Equal(t, 25, res.Inner.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age int `path:"age"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/user/invalid", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/user/:age")
				assert.Equal(t, 0, res.Inner.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Age *int `path:"age"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/user/25", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/user/:age?")
				assert.NotNil(t, res.Age)
				assert.Equal(t, 25, *res.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Age *int `path:"age"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/user/", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/user/:age?")
				assert.Nil(t, res.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Age *int `path:"age"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/user/invalid", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/user/:age?")
				assert.Nil(t, res.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age *int `path:"age"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/user/25", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/user/:age?")
				assert.NotNil(t, res.Inner.Age)
				assert.Equal(t, 25, *res.Inner.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age *int `path:"age"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/user/", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/user/:age?")
				assert.Nil(t, res.Inner.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age *int `path:"age"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/user/invalid", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/user/:age?")
				assert.Nil(t, res.Inner.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
}

func testPathParamOfTypeUint(t *testing.T, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Count uint `path:"count"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/item/10", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/item/:count")
				assert.Equal(t, uint(10), res.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Count uint `path:"count"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/item/invalid", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/item/:count")
				assert.Equal(t, uint(0), res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("negative_value", func(t *testing.T) {
				type Input struct {
					Count uint `path:"count"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/item/-5", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/item/:count")
				assert.Equal(t, uint(0), res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count uint `path:"count"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/item/10", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/item/:count")
				assert.Equal(t, uint(10), res.Inner.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count uint `path:"count"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/item/invalid", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/item/:count")
				assert.Equal(t, uint(0), res.Inner.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Count *uint `path:"count"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/item/10", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/item/:count?")
				assert.NotNil(t, res.Count)
				assert.Equal(t, uint(10), *res.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Count *uint `path:"count"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/item/", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/item/:count?")
				assert.Nil(t, res.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Count *uint `path:"count"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/item/invalid", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/item/:count?")
				assert.Nil(t, res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count *uint `path:"count"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/item/10", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/item/:count?")
				assert.NotNil(t, res.Inner.Count)
				assert.Equal(t, uint(10), *res.Inner.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count *uint `path:"count"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/item/", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/item/:count?")
				assert.Nil(t, res.Inner.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count *uint `path:"count"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/item/invalid", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/item/:count?")
				assert.Nil(t, res.Inner.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
}

func testPathParamOfTypeFloat(t *testing.T, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Price float64 `path:"price"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/product/19.99", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/product/:price")
				assert.Equal(t, 19.99, res.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Price float64 `path:"price"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/product/invalid", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/product/:price")
				assert.Equal(t, 0.0, res.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price float64 `path:"price"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/product/19.99", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/product/:price")
				assert.Equal(t, 19.99, res.Inner.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price float64 `path:"price"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/product/invalid", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/product/:price")
				assert.Equal(t, 0.0, res.Inner.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Price *float64 `path:"price"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/product/19.99", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/product/:price?")
				assert.NotNil(t, res.Price)
				assert.Equal(t, 19.99, *res.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Price *float64 `path:"price"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/product/", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/product/:price?")
				assert.Nil(t, res.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Price *float64 `path:"price"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/product/invalid", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/product/:price?")
				assert.Nil(t, res.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price *float64 `path:"price"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/product/19.99", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/product/:price?")
				assert.NotNil(t, res.Inner.Price)
				assert.Equal(t, 19.99, *res.Inner.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price *float64 `path:"price"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/product/", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/product/:price?")
				assert.Nil(t, res.Inner.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price *float64 `path:"price"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/product/invalid", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/product/:price?")
				assert.Nil(t, res.Inner.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
}

func testPathParamOfTypeBool(t *testing.T, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value_true", func(t *testing.T) {
				type Input struct {
					Active bool `path:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/status/true", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/status/:active")
				assert.Equal(t, true, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("with_value_false", func(t *testing.T) {
				type Input struct {
					Active bool `path:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/status/false", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/status/:active")
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("with_value_1", func(t *testing.T) {
				type Input struct {
					Active bool `path:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/status/1", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/status/:active")
				assert.Equal(t, true, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("with_value_0", func(t *testing.T) {
				type Input struct {
					Active bool `path:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/status/0", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/status/:active")
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Active bool `path:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/status/invalid", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/status/:active")
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active bool `path:"active"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/status/true", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/status/:active")
				assert.Equal(t, true, res.Inner.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active bool `path:"active"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/status/invalid", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/status/:active")
				assert.Equal(t, false, res.Inner.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value_true", func(t *testing.T) {
				type Input struct {
					Active *bool `path:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/status/true", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/status/:active?")
				assert.NotNil(t, res.Active)
				assert.Equal(t, true, *res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("with_value_false", func(t *testing.T) {
				type Input struct {
					Active *bool `path:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/status/false", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/status/:active?")
				assert.NotNil(t, res.Active)
				assert.Equal(t, false, *res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Active *bool `path:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/status/", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/status/:active?")
				assert.Nil(t, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Active *bool `path:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/status/invalid", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/status/:active?")
				assert.Nil(t, res.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active *bool `path:"active"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/status/true", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/status/:active?")
				assert.NotNil(t, res.Inner.Active)
				assert.Equal(t, true, *res.Inner.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active *bool `path:"active"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/status/", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/status/:active?")
				assert.Nil(t, res.Inner.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active *bool `path:"active"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/status/invalid", nil)
				res, status, _ := doPathTest[Input](t, app, req, routeRegistration, "/status/:active?")
				assert.Nil(t, res.Inner.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
}
