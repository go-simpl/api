package injectiontests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-simpl/simplapi/tests/utils"
)

func TestCookieParam(t *testing.T) {
	frameworks := utils.GetAllFrameworks()
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}

	types := utils.GetSupportedTypes()

	for _, framework := range frameworks {
		t.Run(framework, func(t *testing.T) {
			for _, method := range methods {
				t.Run(method, func(t *testing.T) {
					for _, typeName := range types {
						t.Run(typeName, func(t *testing.T) {
							getCookieParamTestForType(typeName)(t, framework, method)
						})
					}
				})
			}
		})
	}
}

func getCookieParamTestForType(typeName string) func(t *testing.T, frameworkName string, method string) {
	switch typeName {
	case "string":
		return testCookieParamOfTypeString
	case "int":
		return testCookieParamOfTypeInt
	case "uint":
		return testCookieParamOfTypeUint
	case "float64":
		return testCookieParamOfTypeFloat
	case "bool":
		return testCookieParamOfTypeBool
	default:
		panic("tests not defined")
	}
}

func testCookieParamOfTypeString(t *testing.T, frameworkName string, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Name string `cookie:"X-Name"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Name", Value: "John"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, "John", res.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Name string `cookie:"X-Name"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, "", res.Name)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Name string `cookie:"X-Name"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Name", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, "", res.Name)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name string `cookie:"X-Name"`
					}
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Name", Value: "John"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, "John", res.Inner.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name string `cookie:"X-Name"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, "", res.Inner.Name)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name string `cookie:"X-Name"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Name", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, "", res.Inner.Name)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Name *string `cookie:"X-Name"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Name", Value: "John"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.NotNil(t, res.Name)
				assert.Equal(t, "John", *res.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Name *string `cookie:"X-Name"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Name *string `cookie:"X-Name"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Name", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Name)
				assert.Equal(t, http.StatusOK, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name *string `cookie:"X-Name"`
					}
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Name", Value: "John"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.NotNil(t, res.Inner.Name)
				assert.Equal(t, "John", *res.Inner.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name *string `cookie:"X-Name"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Inner.Name)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Name *string `cookie:"X-Name"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Name", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Inner.Name)
				assert.Equal(t, http.StatusOK, status)
			})
		})
	})
}

func testCookieParamOfTypeInt(t *testing.T, frameworkName string, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Age int `cookie:"X-Age"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Age", Value: "25"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, 25, res.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Age int `cookie:"X-Age"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, 0, res.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Age int `cookie:"X-Age"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Age", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, 0, res.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Age int `cookie:"X-Age"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Age", Value: "invalid"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, 0, res.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age int `cookie:"X-Age"`
					}
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Age", Value: "25"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, 25, res.Inner.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age int `cookie:"X-Age"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, 0, res.Inner.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age int `cookie:"X-Age"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Age", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, 0, res.Inner.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age int `cookie:"X-Age"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Age", Value: "invalid"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, 0, res.Inner.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Age *int `cookie:"X-Age"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Age", Value: "25"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.NotNil(t, res.Age)
				assert.Equal(t, 25, *res.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Age *int `cookie:"X-Age"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Age *int `cookie:"X-Age"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Age", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Age *int `cookie:"X-Age"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Age", Value: "invalid"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Age)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age *int `cookie:"X-Age"`
					}
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Age", Value: "25"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.NotNil(t, res.Inner.Age)
				assert.Equal(t, 25, *res.Inner.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age *int `cookie:"X-Age"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Inner.Age)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Age *int `cookie:"X-Age"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Age", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Inner.Age)
				assert.Equal(t, http.StatusOK, status)
			})
		})
	})
}

func testCookieParamOfTypeUint(t *testing.T, frameworkName string, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Count uint `cookie:"X-Count"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Count", Value: "10"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, uint(10), res.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Count uint `cookie:"X-Count"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, uint(0), res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Count uint `cookie:"X-Count"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Count", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, uint(0), res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Count uint `cookie:"X-Count"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Count", Value: "invalid"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, uint(0), res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("negative_value", func(t *testing.T) {
				type Input struct {
					Count uint `cookie:"X-Count"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Count", Value: "-5"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, uint(0), res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count uint `cookie:"X-Count"`
					}
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Count", Value: "10"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, uint(10), res.Inner.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count uint `cookie:"X-Count"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, uint(0), res.Inner.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count uint `cookie:"X-Count"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Count", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, uint(0), res.Inner.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count uint `cookie:"X-Count"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Count", Value: "invalid"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, uint(0), res.Inner.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Count *uint `cookie:"X-Count"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Count", Value: "10"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.NotNil(t, res.Count)
				assert.Equal(t, uint(10), *res.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Count *uint `cookie:"X-Count"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Count *uint `cookie:"X-Count"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Count", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Count *uint `cookie:"X-Count"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Count", Value: "invalid"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count *uint `cookie:"X-Count"`
					}
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Count", Value: "10"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.NotNil(t, res.Inner.Count)
				assert.Equal(t, uint(10), *res.Inner.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count *uint `cookie:"X-Count"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Inner.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count *uint `cookie:"X-Count"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Count", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Inner.Count)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Count *uint `cookie:"X-Count"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Count", Value: "invalid"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Inner.Count)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
}

func testCookieParamOfTypeFloat(t *testing.T, frameworkName string, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Price float64 `cookie:"X-Price"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Price", Value: "19.99"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, 19.99, res.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Price float64 `cookie:"X-Price"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, 0.0, res.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Price float64 `cookie:"X-Price"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Price", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, 0.0, res.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Price float64 `cookie:"X-Price"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Price", Value: "invalid"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, 0.0, res.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price float64 `cookie:"X-Price"`
					}
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Price", Value: "19.99"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, 19.99, res.Inner.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price float64 `cookie:"X-Price"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, 0.0, res.Inner.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price float64 `cookie:"X-Price"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Price", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, 0.0, res.Inner.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price float64 `cookie:"X-Price"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Price", Value: "invalid"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, 0.0, res.Inner.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Price *float64 `cookie:"X-Price"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Price", Value: "19.99"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.NotNil(t, res.Price)
				assert.Equal(t, 19.99, *res.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Price *float64 `cookie:"X-Price"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Price *float64 `cookie:"X-Price"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Price", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Price *float64 `cookie:"X-Price"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Price", Value: "invalid"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price *float64 `cookie:"X-Price"`
					}
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Price", Value: "19.99"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.NotNil(t, res.Inner.Price)
				assert.Equal(t, 19.99, *res.Inner.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price *float64 `cookie:"X-Price"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Inner.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price *float64 `cookie:"X-Price"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Price", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Inner.Price)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Price *float64 `cookie:"X-Price"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Price", Value: "invalid"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Inner.Price)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
}

func testCookieParamOfTypeBool(t *testing.T, frameworkName string, method string) {
	t.Run("required", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value_true", func(t *testing.T) {
				type Input struct {
					Active bool `cookie:"X-Active"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Active", Value: "true"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, true, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("with_value_false", func(t *testing.T) {
				type Input struct {
					Active bool `cookie:"X-Active"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Active", Value: "false"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("with_value_1", func(t *testing.T) {
				type Input struct {
					Active bool `cookie:"X-Active"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Active", Value: "1"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, true, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("with_value_0", func(t *testing.T) {
				type Input struct {
					Active bool `cookie:"X-Active"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Active", Value: "0"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Active bool `cookie:"X-Active"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Active bool `cookie:"X-Active"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Active", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Active bool `cookie:"X-Active"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Active", Value: "invalid"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, false, res.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active bool `cookie:"X-Active"`
					}
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Active", Value: "true"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, true, res.Inner.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active bool `cookie:"X-Active"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, false, res.Inner.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active bool `cookie:"X-Active"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Active", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, false, res.Inner.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active bool `cookie:"X-Active"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Active", Value: "invalid"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Equal(t, false, res.Inner.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
	t.Run("optional", func(t *testing.T) {
		t.Run("direct_struct", func(t *testing.T) {
			t.Run("with_value_true", func(t *testing.T) {
				type Input struct {
					Active *bool `cookie:"X-Active"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Active", Value: "true"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.NotNil(t, res.Active)
				assert.Equal(t, true, *res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("with_value_false", func(t *testing.T) {
				type Input struct {
					Active *bool `cookie:"X-Active"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Active", Value: "false"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.NotNil(t, res.Active)
				assert.Equal(t, false, *res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Active *bool `cookie:"X-Active"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Active *bool `cookie:"X-Active"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Active", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Active *bool `cookie:"X-Active"`
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Active", Value: "invalid"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
		t.Run("nested_struct", func(t *testing.T) {
			t.Run("with_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active *bool `cookie:"X-Active"`
					}
				}

				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Active", Value: "true"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.NotNil(t, res.Inner.Active)
				assert.Equal(t, true, *res.Inner.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("without_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active *bool `cookie:"X-Active"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Inner.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("empty_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active *bool `cookie:"X-Active"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Active", Value: ""})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Inner.Active)
				assert.Equal(t, http.StatusOK, status)
			})
			t.Run("invalid_value", func(t *testing.T) {
				type Input struct {
					Inner struct {
						Active *bool `cookie:"X-Active"`
					}
				}
				app := SetupTest(frameworkName, method)
				req := httptest.NewRequest(method, "/", nil)
				req.AddCookie(&http.Cookie{Name: "X-Active", Value: "invalid"})
				res, status, _ := doTest[Input](t, app, req, "/")
				assert.Nil(t, res.Inner.Active)
				assert.Equal(t, http.StatusUnprocessableEntity, status)
			})
		})
	})
}
