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

func TestQueryParam(t *testing.T) {
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
					t.Run("string", func(t *testing.T) { testQueryParamOfTypeString(t, method) })
					t.Run("int", func(t *testing.T) { testQueryParamOfTypeInt(t, method) })
					t.Run("uint", func(t *testing.T) { testQueryParamOfTypeUint(t, method) })
					t.Run("float", func(t *testing.T) { testQueryParamOfTypeFloat(t, method) })
					t.Run("bool", func(t *testing.T) { testQueryParamOfTypeBool(t, method) })
				})
			}
		})
	}
}

func doTest[T any](t *testing.T, app *simplapi.App, req *http.Request, routeRegistration RouteRegistrationFunc) (T, int, string) {
	var result T

	routeRegistration(req.URL.Path, nil, func(input T) (*HelloResponse, error) {
		result = input
		return &HelloResponse{Message: "Hello"}, nil
	})

	resp, err := app.GetApp().TestRequest(req)
	assert.NoError(t, err)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	return result, resp.StatusCode, string(bodyBytes)
}

func testQueryParamOfTypeString(t *testing.T, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Name string `query:"name"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?name=John", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, "John", res.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Name string `query:"name"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, "", res.Name)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Name string `query:"name"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/name=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, "", res.Name)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name string `query:"name"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?name=John", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, "John", res.Inner.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name string `query:"name"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, "", res.Inner.Name)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name string `query:"name"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/name=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, "", res.Inner.Name)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Name *string `query:"name"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?name=John", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.NotNil(t, res.Name)
				assert.Equal(t, "John", *res.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Name *string `query:"name"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Name *string `query:"name"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/name=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Name)
				assert.Equal(t, http.StatusOK, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name *string `query:"name"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?name=John", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.NotNil(t, res.Inner.Name)
				assert.Equal(t, "John", *res.Inner.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name *string `query:"name"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Inner.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name *string `query:"name"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/name=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Inner.Name)
				assert.Equal(t, http.StatusOK, status)
			})
		})
	})
}

func testQueryParamOfTypeInt(t *testing.T, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Age int `query:"age"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?age=25", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, 25, res.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Age int `query:"age"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, 0, res.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Age int `query:"age"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?age=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, 0, res.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Age int `query:"age"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?age=invalid", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, 0, res.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age int `query:"age"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?age=25", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, 25, res.Inner.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age int `query:"age"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, 0, res.Inner.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age int `query:"age"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?age=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, 0, res.Inner.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age int `query:"age"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?age=invalid", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, 0, res.Inner.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Age *int `query:"age"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?age=25", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.NotNil(t, res.Age)
				assert.Equal(t, 25, *res.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Age *int `query:"age"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Age *int `query:"age"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?age=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Age *int `query:"age"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?age=invalid", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age *int `query:"age"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?age=25", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.NotNil(t, res.Inner.Age)
				assert.Equal(t, 25, *res.Inner.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age *int `query:"age"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Inner.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age *int `query:"age"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?age=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Inner.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age *int `query:"age"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?age=invalid", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Inner.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
}

func testQueryParamOfTypeUint(t *testing.T, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Count uint `query:"count"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?count=10", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, uint(10), res.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Count uint `query:"count"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, uint(0), res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Count uint `query:"count"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?count=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, uint(0), res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Count uint `query:"count"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?count=invalid", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, uint(0), res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("negative_value", func(t *testing.T) {
				type Input struct {
					Count uint `query:"count"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?count=-5", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, uint(0), res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count uint `query:"count"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?count=10", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, uint(10), res.Inner.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count uint `query:"count"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, uint(0), res.Inner.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count uint `query:"count"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?count=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, uint(0), res.Inner.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count uint `query:"count"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?count=invalid", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, uint(0), res.Inner.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Count *uint `query:"count"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?count=10", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.NotNil(t, res.Count)
				assert.Equal(t, uint(10), *res.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Count *uint `query:"count"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Count *uint `query:"count"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?count=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Count *uint `query:"count"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?count=invalid", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count *uint `query:"count"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?count=10", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.NotNil(t, res.Inner.Count)
				assert.Equal(t, uint(10), *res.Inner.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count *uint `query:"count"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Inner.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count *uint `query:"count"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?count=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Inner.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count *uint `query:"count"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?count=invalid", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Inner.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
}

func testQueryParamOfTypeFloat(t *testing.T, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Price float64 `query:"price"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?price=19.99", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, 19.99, res.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Price float64 `query:"price"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, 0.0, res.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Price float64 `query:"price"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?price=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, 0.0, res.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Price float64 `query:"price"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?price=invalid", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, 0.0, res.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price float64 `query:"price"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?price=19.99", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, 19.99, res.Inner.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price float64 `query:"price"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, 0.0, res.Inner.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price float64 `query:"price"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?price=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, 0.0, res.Inner.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price float64 `query:"price"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?price=invalid", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, 0.0, res.Inner.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Price *float64 `query:"price"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?price=19.99", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.NotNil(t, res.Price)
				assert.Equal(t, 19.99, *res.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Price *float64 `query:"price"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Price *float64 `query:"price"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?price=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Price *float64 `query:"price"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?price=invalid", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price *float64 `query:"price"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?price=19.99", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.NotNil(t, res.Inner.Price)
				assert.Equal(t, 19.99, *res.Inner.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price *float64 `query:"price"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Inner.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price *float64 `query:"price"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?price=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Inner.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price *float64 `query:"price"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?price=invalid", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Inner.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
}

func testQueryParamOfTypeBool(t *testing.T, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value_true", func(t *testing.T) {
				type Input struct {
					Active bool `query:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?active=true", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, true, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("with_value_false", func(t *testing.T) {
				type Input struct {
					Active bool `query:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?active=false", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("with_value_1", func(t *testing.T) {
				type Input struct {
					Active bool `query:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?active=1", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, true, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("with_value_0", func(t *testing.T) {
				type Input struct {
					Active bool `query:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?active=0", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Active bool `query:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Active bool `query:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?active=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Active bool `query:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?active=invalid", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active bool `query:"active"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?active=true", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, true, res.Inner.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active bool `query:"active"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, false, res.Inner.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active bool `query:"active"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?active=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, false, res.Inner.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active bool `query:"active"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?active=invalid", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Equal(t, false, res.Inner.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value_true", func(t *testing.T) {
				type Input struct {
					Active *bool `query:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?active=true", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.NotNil(t, res.Active)
				assert.Equal(t, true, *res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("with_value_false", func(t *testing.T) {
				type Input struct {
					Active *bool `query:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?active=false", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.NotNil(t, res.Active)
				assert.Equal(t, false, *res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Active *bool `query:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Active *bool `query:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?active=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Active *bool `query:"active"`
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?active=invalid", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active *bool `query:"active"`
					}
				}

				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?active=true", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.NotNil(t, res.Inner.Active)
				assert.Equal(t, true, *res.Inner.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active *bool `query:"active"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Inner.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active *bool `query:"active"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?active=", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Inner.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active *bool `query:"active"`
					}
				}
				app, routeRegistration := SetupTest(method)
				req := httptest.NewRequest(method, "/?active=invalid", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration)
				assert.Nil(t, res.Inner.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
}
