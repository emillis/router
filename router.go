package router

import (
	"errors"
	"fmt"
	"net/http"
)

type Segment struct {
	value      string
	isVariable bool
	ok         bool
}

type Route struct {
	original     string
	segments     []Segment
	hasVariables bool
}

func (r *Route) Compare(path string) bool {
	j := 0
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] != 47 {
			continue
		}

		if !r.segments[j].isVariable && r.segments[j].value != path[i:] {
			fmt.Println(1)
			return false
		}

		path = path[:i]
		j++
		if j >= len(r.segments) {
			fmt.Println(2)
			return false
		}
	}
	if j != len(r.segments) {
		fmt.Println(3)
		return false
	}
	return true
}

type Router struct {
	routes []*Route
}

func (r *Router) findRoute(s string) *Route {
	s, _ = processPath(s)

	for i := 0; i < len(r.routes); i++ {
		if !r.routes[i].Compare(s) {
			continue
		}

		return r.routes[i]
	}

	return nil
}

func (r *Router) addRoute(s string) error {
	route, err := NewRoute(s)
	if err != nil {
		return err
	}

	r.routes = append(r.routes, route)

	return nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("asda"))
	fmt.Println(req.URL)
	fmt.Println(req.URL.Path)
	fmt.Println(req.URL.RawQuery)
	fmt.Println(req.URL.Query())
	fmt.Println(req.URL.RequestURI())
	fmt.Println(req.URL.String())
}

//===========[FUNCTIONALITY]====================================================================================================

func NewSegment(segment string) Segment {
	//TODO fix segment being "". It would cause problems below getting the index 0 of the string

	return Segment{
		value:      segment,
		isVariable: segment[1] == 58,
		ok:         true,
	}
}

func SplitPath(path string) []Segment {
	var buffer []Segment
	var j int

	for i := len(path) - 1; i >= 0; i-- {
		if path[i] != 47 {
			continue
		}

		buffer = append(buffer, NewSegment(path[i:]))
		path = path[:i]
		i = len(path)

		j++
	}

	return buffer
}

//processPath check for critical errors within the path supplied. Also, removes trailing "/" sign if present
func processPath(s string) (string, error) {
	if s == "" {
		return s, errors.New("path supplied cannot be an empty string")
	}

	if s[0] != 47 {
		return s, errors.New("path must begin with \"/\"")
	}

	if s[len(s)-1] == 47 && len(s) > 1 {
		s = s[:len(s)-1]
	}

	return s, nil
}

func NewRoute(s string) (*Route, error) {
	s, err := processPath(s)
	if err != nil {
		return nil, err
	}

	r := Route{
		original: s,
		segments: SplitPath(s),
	}

	for _, segment := range r.segments {
		if !segment.isVariable {
			continue
		}

		r.hasVariables = true
		break
	}

	return &r, nil
}

func NewRouter() *Router {
	return &Router{
		routes: []*Route{},
	}
}
