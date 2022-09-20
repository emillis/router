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
	original        string
	partitionedPath []string
	segments        []Segment
	hasVariables    bool
}

func (r *Route) AddSegment(s Segment) {

}

func ParseSegment(segment string) Segment {
	//TODO fix segment being "". It would cause problems below getting the index 0 of the string

	s := Segment{
		value:      segment,
		isVariable: segment[0] == 58,
	}

	return s
}

func ParseRoute(route string) Route {
	r := Route{
		original: route,
	}

	//Stripping the first "/" if exist
	if route[0] == 47 {
		route = route[1:]
	}

	//Stripping the last "/" if exist
	if route[len(route)-1] == 47 {
		route = route[:len(route)-1]
	}

	r.partitionedPath = strings.Split(route, "/")

	return r
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
	s := []Segment{{"1", false}, {"2", false}, {"3", false}}

	for n := 0; n < b.N; n++ {
		for _, v := range s {
			if v.value == "-999" {

			}
		}
	}
}
