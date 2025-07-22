package injectiontests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-simpl/simplapi/tests/utils"
)

func TestHeaderParam(t *testing.T) {
	frameworks := utils.GetAllFrameworks()
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

	types := utils.GetSupportedTypes()

	for _, framework := range frameworks {
		t.Run(framework, func(t *testing.T) {
			for _, method := range methods {
				t.Run(method, func(t *testing.T) {
					for _, typeName := range types {
						t.Run(typeName, func(t *testing.T) {
							getHeaderParamTestForType(typeName)(t, framework, method)
						})
					}
				})
			}
		})
	}
}

func getHeaderParamTestForType(typeName string) func(t *testing.T, frameworkName string, method string) {
	switch typeName {
	case "string":
		return testHeaderParamOfTypeString
	case "int":
		return testHeaderParamOfTypeInt
	case "uint":
		return testHeaderParamOfTypeUint
	case "float64":
		return testHeaderParamOfTypeFloat
	case "bool":
		return testHeaderParamOfTypeBool
	default:
		panic("tests not defined")
	}
}

func testHeaderParamOfTypeString(t *testing.T, frameworkName string, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Name string `header:"X-Name"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Name", "John")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, "John", res.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Name string `header:"X-Name"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, "", res.Name)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Name string `header:"X-Name"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Name", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, "", res.Name)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name string `header:"X-Name"`
					}
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Name", "John")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, "John", res.Inner.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name string `header:"X-Name"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, "", res.Inner.Name)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name string `header:"X-Name"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Name", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, "", res.Inner.Name)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Name *string `header:"X-Name"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Name", "John")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.NotNil(t, res.Name)
				assert.Equal(t, "John", *res.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Name *string `header:"X-Name"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Name *string `header:"X-Name"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Name", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Name)
				assert.Equal(t, http.StatusOK, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name *string `header:"X-Name"`
					}
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Name", "John")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.NotNil(t, res.Inner.Name)
				assert.Equal(t, "John", *res.Inner.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name *string `header:"X-Name"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Inner.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name *string `header:"X-Name"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Name", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Inner.Name)
				assert.Equal(t, http.StatusOK, status)
			})
		})
	})
}

func testHeaderParamOfTypeInt(t *testing.T, frameworkName string, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Age int `header:"X-Age"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Age", "25")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, 25, res.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Age int `header:"X-Age"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, 0, res.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Age int `header:"X-Age"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Age", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, 0, res.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Age int `header:"X-Age"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Age", "invalid")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, 0, res.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age int `header:"X-Age"`
					}
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Age", "25")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, 25, res.Inner.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age int `header:"X-Age"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, 0, res.Inner.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age int `header:"X-Age"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Age", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, 0, res.Inner.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age int `header:"X-Age"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Age", "invalid")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, 0, res.Inner.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Age *int `header:"X-Age"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Age", "25")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.NotNil(t, res.Age)
				assert.Equal(t, 25, *res.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Age *int `header:"X-Age"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Age *int `header:"X-Age"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Age", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Age *int `header:"X-Age"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Age", "invalid")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age *int `header:"X-Age"`
					}
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Age", "25")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.NotNil(t, res.Inner.Age)
				assert.Equal(t, 25, *res.Inner.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age *int `header:"X-Age"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Inner.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age *int `header:"X-Age"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Age", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Inner.Age)
				assert.Equal(t, http.StatusOK, status)
			})
		})
	})
}

func testHeaderParamOfTypeUint(t *testing.T, frameworkName string, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Count uint `header:"X-Count"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Count", "10")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, uint(10), res.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Count uint `header:"X-Count"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, uint(0), res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Count uint `header:"X-Count"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Count", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, uint(0), res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Count uint `header:"X-Count"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Count", "invalid")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, uint(0), res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("negative_value", func(t *testing.T) {
				type Input struct {
					Count uint `header:"X-Count"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Count", "-5")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, uint(0), res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count uint `header:"X-Count"`
					}
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Count", "10")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, uint(10), res.Inner.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count uint `header:"X-Count"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, uint(0), res.Inner.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count uint `header:"X-Count"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Count", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, uint(0), res.Inner.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count uint `header:"X-Count"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Count", "invalid")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, uint(0), res.Inner.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Count *uint `header:"X-Count"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Count", "10")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.NotNil(t, res.Count)
				assert.Equal(t, uint(10), *res.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Count *uint `header:"X-Count"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Count *uint `header:"X-Count"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Count", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Count *uint `header:"X-Count"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Count", "invalid")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count *uint `header:"X-Count"`
					}
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Count", "10")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.NotNil(t, res.Inner.Count)
				assert.Equal(t, uint(10), *res.Inner.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count *uint `header:"X-Count"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Inner.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count *uint `header:"X-Count"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Count", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Inner.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count *uint `header:"X-Count"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Count", "invalid")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Inner.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
}

func testHeaderParamOfTypeFloat(t *testing.T, frameworkName string, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Price float64 `header:"X-Price"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Price", "19.99")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, 19.99, res.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Price float64 `header:"X-Price"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, 0.0, res.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Price float64 `header:"X-Price"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Price", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, 0.0, res.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Price float64 `header:"X-Price"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Price", "invalid")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, 0.0, res.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price float64 `header:"X-Price"`
					}
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Price", "19.99")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, 19.99, res.Inner.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price float64 `header:"X-Price"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, 0.0, res.Inner.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price float64 `header:"X-Price"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Price", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, 0.0, res.Inner.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price float64 `header:"X-Price"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Price", "invalid")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, 0.0, res.Inner.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Price *float64 `header:"X-Price"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Price", "19.99")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.NotNil(t, res.Price)
				assert.Equal(t, 19.99, *res.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Price *float64 `header:"X-Price"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Price *float64 `header:"X-Price"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Price", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Price *float64 `header:"X-Price"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Price", "invalid")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price *float64 `header:"X-Price"`
					}
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Price", "19.99")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.NotNil(t, res.Inner.Price)
				assert.Equal(t, 19.99, *res.Inner.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price *float64 `header:"X-Price"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Inner.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price *float64 `header:"X-Price"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Price", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Inner.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price *float64 `header:"X-Price"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Price", "invalid")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Inner.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
}

func testHeaderParamOfTypeBool(t *testing.T, frameworkName string, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value_true", func(t *testing.T) {
				type Input struct {
					Active bool `header:"X-Active"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Active", "true")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, true, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("with_value_false", func(t *testing.T) {
				type Input struct {
					Active bool `header:"X-Active"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Active", "false")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("with_value_1", func(t *testing.T) {
				type Input struct {
					Active bool `header:"X-Active"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Active", "1")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, true, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("with_value_0", func(t *testing.T) {
				type Input struct {
					Active bool `header:"X-Active"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Active", "0")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Active bool `header:"X-Active"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Active bool `header:"X-Active"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Active", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Active bool `header:"X-Active"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Active", "invalid")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active bool `header:"X-Active"`
					}
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Active", "true")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, true, res.Inner.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active bool `header:"X-Active"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, false, res.Inner.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active bool `header:"X-Active"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Active", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, false, res.Inner.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active bool `header:"X-Active"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Active", "invalid")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Equal(t, false, res.Inner.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value_true", func(t *testing.T) {
				type Input struct {
					Active *bool `header:"X-Active"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Active", "true")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.NotNil(t, res.Active)
				assert.Equal(t, true, *res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("with_value_false", func(t *testing.T) {
				type Input struct {
					Active *bool `header:"X-Active"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Active", "false")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.NotNil(t, res.Active)
				assert.Equal(t, false, *res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Active *bool `header:"X-Active"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Active *bool `header:"X-Active"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Active", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Active *bool `header:"X-Active"`
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Active", "invalid")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active *bool `header:"X-Active"`
					}
				}

				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Active", "true")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.NotNil(t, res.Inner.Active)
				assert.Equal(t, true, *res.Inner.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active *bool `header:"X-Active"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Inner.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active *bool `header:"X-Active"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Active", "")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Inner.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active *bool `header:"X-Active"`
					}
				}
				app, routeRegistration := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.Header.Set("X-Active", "invalid")
				res, status, _ := doTest[Input](t, app, req, routeRegistration, "/")
				assert.Nil(t, res.Inner.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
}
