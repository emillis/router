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

type Route struct {
	partitionedPath []string
}

func New(path string) *Route {
	tidyString := path

	if tidyString[0] == 47 {
		tidyString = tidyString[1:]
	}

	if tidyString[len(tidyString)-1] == 47 {
		tidyString = tidyString[:len(tidyString)-1]
	}

	return &Route{
		partitionedPath: strings.Split(tidyString, "/"),
	}
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
