package veryFastRouter

import (
	"net/http"
	"net/url"
	"testing"
)

func BenchmarkHttpRouter_ServeHTTP(b *testing.B) {
	router := NewRouter()

	router.HandleFunc("/test", []string{http.MethodGet}, func(w http.ResponseWriter, r *http.Request, info *AdditionalInfo) {})
	router.HandleFunc("/test/two", []string{http.MethodGet}, func(w http.ResponseWriter, r *http.Request, info *AdditionalInfo) {})
	router.HandleFunc("/test/:two/three", []string{http.MethodGet}, func(w http.ResponseWriter, r *http.Request, info *AdditionalInfo) {})
	router.HandleFunc("/test/two/three/four", []string{http.MethodGet}, func(w http.ResponseWriter, r *http.Request, info *AdditionalInfo) {})

	mockRequest := http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/test/two"},
	}

	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		router.ServeHTTP(nil, &mockRequest)
	}
}

func BenchmarkProcessPath(b *testing.B) {
	path := "/one/two/three/four/"

	for n := 0; n < b.N; n++ {

		removeTrailingSlash(path)

	}
}

func BenchmarkRouter_findRoute(b *testing.B) {
	router := NewRouter()
	routesToAdd := []string{
		"/one/:two/three",
		"/one/two/three/four/",
		"/one/two/three/four/five/",
	}
	for _, r := range routesToAdd {
		_, err := router.addRoute(r)
		if err != nil {
			panic(err)
		}
	}

	path := "/one/two/three/four"

	for n := 0; n < b.N; n++ {
		router.findRoute(path)
	}
}

func TestHttpRouter_AddingHandlers(t *testing.T) {
	r := NewRouter()
	requiredCountStatic := 1
	requiredCountVariable := 2

	r.HandleFunc("/one", []string{"GET", "POST"}, func(w http.ResponseWriter, r *http.Request, info *AdditionalInfo) {})
	r.HandleFunc("/one/:two", []string{"GET", "POST"}, func(w http.ResponseWriter, r *http.Request, info *AdditionalInfo) {})
	r.HandleFunc("/one/two/*three", []string{"GET", "POST"}, func(w http.ResponseWriter, r *http.Request, info *AdditionalInfo) {})

	if len(r.variableRoutes) != requiredCountVariable {
		t.Errorf("Number of variable routes should be %d, got %d", requiredCountVariable, len(r.variableRoutes))
	}

	if len(r.staticRoutes) != requiredCountStatic {
		t.Errorf("Number of static routes should be %d, got %d", requiredCountStatic, len(r.staticRoutes))
	}
}

func TestHttpRouter_Routing(t *testing.T) {
	r := NewRouter()

	addHandleFunc := map[string]int{
		"/one/two/":      3,
		"/one/two/three": 6,
		"/one":           11,
		"/":              7,

		"/:one":                  10,
		"/one/:two/three/":       5,
		"/two/:three/:four/":     12,
		"/three/:four/five/:six": 2,
		"/four/five/:six/":       8,
		"/five/:six/:seven/":     9,

		"/one/two/*three": 4,
		"/two/*three":     1,
	}

	//-1 Means pattern not found
	tests := map[string]int{
		"/one/two":         3,
		"/":                7,
		"/one":             11,
		"/one/two/three/":  6,
		"/one/one/one/one": -1,

		"/hello":             10,
		"/one/sixteen/three": 5,
		"/two/nine/nine":     12,

		"/one/two/three/four/five/six/": 4,
		"/one/none/":                    -1,
	}

	res := 0

	newHandleFunc := func(router *HttpRouter, pattern string, val int) {
		router.HandleFunc(pattern, []string{"GET"}, func(w http.ResponseWriter, r *http.Request, info *AdditionalInfo) {
			res = val
		})
	}

	for key, val := range addHandleFunc {
		newHandleFunc(r, key, val)
	}

	for pattern, val := range tests {
		res = -1
		route, _ := r.findRoute(pattern)
		originalPattern := "Not Found"
		if route != nil {
			route.handler(nil, nil, nil)
			originalPattern = route.originalPattern
		}

		if res != val {
			t.Errorf("Expected result %d, got %d. %s, %s", val, res, originalPattern, pattern)
		}
	}
}
