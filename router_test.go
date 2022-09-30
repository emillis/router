package router

import (
	"fmt"
	"testing"
)

var xxx int

func BenchmarkSplitPath(b *testing.B) {
	path := "/one/:two/three/four"

	for n := 0; n < b.N; n++ {
		SplitPath(path)
	}
}

func BenchmarkNewRoute(b *testing.B) {
	path := "/one/two/three/four/"

	for n := 0; n < b.N; n++ {
		NewRoute(path)
	}
}

//The best so far
func BenchmarkSplitType1(b *testing.B) {
	path := "/one/:two/three/four/"

	//x := HttpRouter{}
	//
	//log.Fatal(http.ListenAndServe(":80", x))

	var segments [4]Segment

	for n := 0; n < b.N; n++ {
		//startPos := 0
		j := 0

		for i := 1; i < len(path); i++ {
			if path[i] != 47 {
				continue
			}

			segments[j] = Segment{
				value:      path[:i],
				isVariable: path[:i][1] == 58,
				ok:         true,
			}
			//segments[j] = NewSegment(path[:i])
			path = path[i:]
			i = 0

			j++
		}
		path = "/one/:two/three/four/"
	}

	fmt.Println(segments)
}

//Same as type 1, but in reverse
func BenchmarkSplitType2(b *testing.B) {

	for n := 0; n < b.N; n++ {
		path := "/one/:two/three/four"
		var segments [bufferSize]Segment

		j := 0

		for i := len(path) - 1; i >= 0; i-- {
			if path[i] != 47 {
				continue
			}

			segments[j] = NewSegment(path[i:])
			path = path[:i]
			i = len(path)

			j++
		}
	}

	//fmt.Println(segments)
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
		err := router.addRoute(r)
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
