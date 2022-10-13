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
