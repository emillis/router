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
		"/*one":              1,
		"/:one/two/three":    2,
		"/one/two/":          3,
		"/one/two/*three":    4,
		"/:one/:two/:three/": 5,
	}

	tests := map[string]int{
		"/hello": 1,
	}

	res := 0

	for key, val := range addHandleFunc {
		r.HandleFunc(key, []string{"GET"}, func(w http.ResponseWriter, req *http.Request, info *AdditionalInfo) {
			res = val
		})
	}

	for key, val := range tests {
		r.ServeHTTP(nil, &http.Request{URL: &url.URL{Path: key}, Method: "GET"})

		if res != val {
			t.Errorf("Expected result %d, got %d", val, res)
		}
	}
}
