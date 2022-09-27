package router

import (
	"fmt"
	"testing"
)

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

	//x := Router{}
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

	}

	fmt.Println(segments)
}

//Same as type 1, but in reverse
func BenchmarkSplitType2(b *testing.B) {
	path, _ := processPath("/one/:two/three/four")

	for n := 0; n < b.N; n++ {
		const bufferSize = 10
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

	//processedString, err := processPath(path)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(processedString)

	for n := 0; n < b.N; n++ {

		processPath(path)

	}
}

func BenchmarkRouter_findRoute(b *testing.B) {
	router := Router{}
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

	//for n := 0; n < b.N; n++ {
	fmt.Println(router.findRoute("/one/two/three/four"))
	//}
}
