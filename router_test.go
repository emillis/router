package veryFastRouter

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

var xxx int

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

func BenchmarkSplitPath(b *testing.B) {
	path := "/one/:two/three/four"

	for n := 0; n < b.N; n++ {
		splitPath(path)
	}
}

func BenchmarkNewRoute(b *testing.B) {
	path := "/one/two/three/four/"

	for n := 0; n < b.N; n++ {
		newRoute(path)
	}
}

func BenchmarkSplitType1(b *testing.B) {
	path := "/one/:two/three/four/"

	//x := HttpRouter{}
	//
	//log.Fatal(http.ListenAndServe(":80", x))

	var segments [4]segment

	for n := 0; n < b.N; n++ {
		//startPos := 0
		j := 0

		for i := 1; i < len(path); i++ {
			if path[i] != 47 {
				continue
			}

			segments[j] = segment{
				original:   path[:i],
				isVariable: path[:i][1] == 58,
			}
			//values[j] = newSegment(path[:i])
			path = path[i:]
			i = 0

			j++
		}
		path = "/one/:two/three/four/"
	}

	fmt.Println(segments)
}

func BenchmarkSplitType2(b *testing.B) {

	for n := 0; n < b.N; n++ {
		path := "/one/:two/three/four"
		var segments [bufferSize]segment

		j := 0

		for i := len(path) - 1; i >= 0; i-- {
			if path[i] != 47 {
				continue
			}

			segments[j] = newSegment(path[i:])
			path = path[:i]
			i = len(path)

			j++
		}
	}

	//fmt.Println(values)
}

func BenchmarkProcessPath(b *testing.B) {
	path := "/one/two/three/four/"

	for n := 0; n < b.N; n++ {

		processPath(path)

	}
}

func BenchmarkReadingFromMap(b *testing.B) {
	var aaa = map[string]bool{
		"/one/two/three/four/":      true,
		"/one/:two/three":           true,
		"/one/two/three/four/five/": true,
	}

	for n := 0; n < b.N; n++ {

		if result, _ := aaa["/one/two/three/four/"]; result {

		}

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

	if xxx == 0 {
		fmt.Println("===============================")
		fmt.Println(router.findRoute(path))
		fmt.Println("===============================")
		xxx++
	}

	for n := 0; n < b.N; n++ {
		router.findRoute(path)
	}
}
