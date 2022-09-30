package router

import (
	"fmt"
	"net/http"
)

//===========[CACHE/STATIC]====================================================================================================

const bufferSize = 50

var GET Method = "GET"
var POST Method = "POST"

//===========[STRUCTS]====================================================================================================

type HandlerFunc func(http.ResponseWriter, *http.Request)
type Method string

type PathDetails struct {
	count    int
	segments [bufferSize]string
}

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

func (r *Route) Compare(pd *PathDetails) bool {
	if pd.count != len(r.segments) {
		return false
	}

	for i := 0; i < pd.count; i++ {
		if pd.segments[i] != r.segments[i].value && !r.segments[i].isVariable {
			return false
		}
	}

	return true
}

type HttpRouter struct {
	staticRoutes   map[string]*Route
	variableRoutes []*Route
}

func (r *HttpRouter) HandleFunc(pattern string, handler HandlerFunc) {

}

func (r *HttpRouter) findRoute(s string) *Route {
	s = processPath(s)

	if router, exist := r.staticRoutes[s]; exist {
		return router
	}

	pd := &PathDetails{
		count:    0,
		segments: [50]string{},
	}

	//Splitting the supplied path into its segments
	for i := len(s) - 1; i >= 0; i-- {
		//If the character is not "/", continue to the next character
		if s[i] != 47 {
			continue
		}

		pd.segments[pd.count] = s[i:]

		s = s[:i]
		pd.count++
	}

	for i := 0; i < len(r.variableRoutes); i++ {
		if !r.variableRoutes[i].Compare(pd) {
			continue
		}

		return r.variableRoutes[i]
	}

	return nil
}

func (r *HttpRouter) addRoute(s string) error {
	route, err := NewRoute(s)
	if err != nil {
		return err
	}

	if !route.hasVariables {
		r.staticRoutes[route.original] = route
		return nil
	}

	r.variableRoutes = append(r.variableRoutes, route)

	return nil
}

func (r *HttpRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
func processPath(s string) string {
	//if s == "" {
	//	return s, errors.New("path supplied cannot be an empty string")
	//}

	//if s[0] != 47 {
	//	return s, errors.New("path must begin with \"/\"")
	//}

	if s[len(s)-1] == 47 && len(s) > 1 {
		return s[:len(s)-1]
	}

	return s
}

func NewRoute(s string) (*Route, error) {
	s = processPath(s)

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

func NewRouter() *HttpRouter {
	return &HttpRouter{
		staticRoutes:   map[string]*Route{},
		variableRoutes: []*Route{},
	}
}
