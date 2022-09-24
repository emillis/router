package router

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

var stringToSplit = "/one/two/three/four/"
var stringSeparator = "/"
var bytesToSplit = []byte("/one/two/three/four")
var bytesSeparator = []byte("/")

type Segment struct {
	value      string
	isVariable bool
}

type Route struct {
	original     string
	segments     []Segment
	hasVariables bool
}

func (r *Route) AddSegment(s Segment) {
	if s.isVariable {
		r.hasVariables = true
	}

	r.segments = append(r.segments, s)
}

func (r *Route) Compare(s []string) {
	//if s.isVariable {
	//	r.hasVariables = true
	//}
	//
	//r.segments = append(r.segments, s)
}

func ParsePath(route string) []string {
	//Stripping the first "/" if exist
	if route[0] == 47 {
		route = route[1:]
	}

	//Stripping the last "/" if exist
	if route[len(route)-1] == 47 {
		route = route[:len(route)-1]
	}

	return strings.Split(route, "/")
}

func CreateSegments(s []string) []Segment {
	var results []Segment

	for _, x := range s {
		results = append(results, NewSegment(x))
	}

	return results
}

func NewRoute(s string) Route {

	r := Route{
		original:     s,
		segments:     []Segment{},
		hasVariables: false,
	}

	for _, segment := range CreateSegments(ParsePath(s)) {
		r.AddSegment(segment)
	}

	return r
}

func NewSegment(segment string) Segment {
	//TODO fix segment being "". It would cause problems below getting the index 0 of the string

	s := Segment{
		value:      segment,
		isVariable: segment[0] == 58,
	}

	return s
}

func BenchmarkEntry_SplitStrings(b *testing.B) {
	tidyString := stringToSplit

	if tidyString[0] == 47 {
		tidyString = tidyString[1:]
	}

	if tidyString[len(tidyString)-1] == 47 {
		tidyString = tidyString[:len(tidyString)-1]
	}

	fmt.Println(tidyString)

	for n := 0; n < b.N; n++ {
		strings.Split(tidyString, stringSeparator)
	}
}

func BenchmarkEntry_SplitBytes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		bytes.Split(bytesToSplit, bytesSeparator)
	}
}

func BenchmarkReadMap(b *testing.B) {
	route := NewRoute("/one")

	for n := 0; n < b.N; n++ {
		for _, v := range route.segments {
			if v.value == "-999" {

			}
		}
	}
}

func BenchmarkSplitType1(b *testing.B) {
	path := "/one/two/three/four"

	for n := 0; n < b.N; n++ {
		var split [4]string
		startPos := 0
		j := 0

		for i, symbol := range path {
			if symbol != 47 {
				continue
			}

			split[j] = path[startPos:i]

			startPos = i + 1
			j++
		}
	}
}
