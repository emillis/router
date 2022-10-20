package veryFastRouter

import (
	"fmt"
	"net/http"
	"sync"
)

//===========[CACHE/STATIC]====================================================================================================

//HandlerFunc defines how a request handler should look like
type HandlerFunc func(http.ResponseWriter, *http.Request, *AdditionalInfo)

//===========[INTERFACES]====================================================================================================

//Authenticator is used to authenticate the request sender's identity
type Authenticator interface {
	Authenticate(*http.Request) bool
}

//Authorizer checks whether the incoming request is authorized to access the route
type Authorizer interface {
	Authorize(*http.Request) bool
}

//===========[STRUCTS]====================================================================================================

//HttpRouter implements Handler interface
type HttpRouter struct {
	//bufferSize is the maximum number of values that the route can consist of. Increasing this doesn't appear
	//to affect performance of the application, only it's memory footprint.
	bufferSize int

	//staticRoutes store all the routes that do not have variables in them
	staticRoutes map[string]*route

	//variableRoutes store all the routes that contain variables or "Match All" pattern in them
	variableRoutes []*route

	//matchAllRoutes store only the routes that have "Match All" pattern
	matchAllRoutes []*route

	//httpStatusCodeHandlers hold all the default/custom handlers to various http status codes
	httpStatusCodeHandlers httpStatusCodeHandlers

	//An empty AdditionalInfo element. Do not use this for storing actual variables
	defaultAdditionalInfo *AdditionalInfo

	//Allows the client to add custom authorization method to the router
	authorizer Authorizer

	//Allows the client to add custom authenticator method in the router
	authenticator Authenticator

	//Mutex used only to add new patterns synchronously
	mx sync.RWMutex
}

//findRoute returns pointer to route based on path supplied as well as a slice of variables
func (r *HttpRouter) findRoute(path string) (*route, []string) {
	path = removeTrailingSlash(path)

	//Looking in static routes first, if there's no match, looks in the routes with variables
	if router, exist := r.staticRoutes[path]; exist {
		return router, nil
	}

	a := make([]string, 0, r.bufferSize)

	//Splitting the supplied path into its values
	for i := len(path) - 1; i >= 0; i-- {
		//If the character is not "/", continue to the next character
		if path[i] != 47 {
			continue
		}

		a = append(a, path[i:])
		path = path[:i]
	}

	//Matching variable routes
	for i := 0; i < len(r.variableRoutes); i++ {
		matched, variables := r.variableRoutes[i].compare(a)

		if !matched {
			continue
		}

		return r.variableRoutes[i], variables
	}

	//Matching "Match All" routes
	for i := 0; i < len(r.matchAllRoutes); i++ {
		matched, variables := r.matchAllRoutes[i].compare(a)

		if !matched {
			continue
		}

		return r.matchAllRoutes[i], variables
	}

	return nil, nil
}

//addRoute parses pattern supplied and adds it to the HttpRouter
func (r *HttpRouter) addRoute(pattern string) (*route, error) {
	route, err := newRoute(pattern)
	if err != nil {
		return nil, err
	}

	if err = r.checkForPathIncongruences(route); err != nil {
		return nil, err
	}

	if route.hasVariables {
		r.variableRoutes = append(r.variableRoutes, route)
		return route, nil
	}

	if route.hasMatchAll {
		r.matchAllRoutes = append(r.matchAllRoutes, route)
		return route, nil
	}

	r.staticRoutes[route.originalPattern] = route
	return route, nil
}

//checkForPathIncongruences checks whether there are no conflicting paths being added
func (r *HttpRouter) checkForPathIncongruences(r2 *route) error {

	if r2.hasVariables {
		for _, r1 := range r.variableRoutes {
			err := r1.compareRoutes(r2)
			if err == nil {
				continue
			}

			return err
		}
		return nil
	}

	if r2.hasMatchAll {
		for _, r1 := range r.matchAllRoutes {
			err := r1.compareRoutes(r2)
			if err == nil {
				continue
			}

			return err
		}
		return nil
	}

	for _, r1 := range r.staticRoutes {
		err := r1.compareRoutes(r2)
		if err == nil {
			continue
		}

		return err
	}

	return nil
}

//AddAuthorizationMethod adds authorization method to the HttpRouter
func (r *HttpRouter) AddAuthorizationMethod(a Authorizer) {
	r.authorizer = a
}

//AddAuthenticationMethod adds authentication method to the HttpRouter
func (r *HttpRouter) AddAuthenticationMethod(a Authenticator) {
	r.authenticator = a
}

//HttpStatusCodeHandler allows you to set up custom handlers for various http status codes, e.g. 404, 405...
func (r *HttpRouter) HttpStatusCodeHandler(statusCode int, handler HandlerFunc) {
	//At first, checking whether the status code exist in the httpStatusCodeHandlers,
	//if not, it means that code is not supported
	if _, exist := r.httpStatusCodeHandlers.handlers[statusCode]; !exist {
		panic(fmt.Sprintf("status code \"%d\" is not supported!", statusCode))
	}

	if handler == nil {
		panic(fmt.Sprintf("handler is not defined for status code \"%d\"!", statusCode))
	}

	//Assigning newly supplied handler in the place of the default one. The purpose of the wrapper
	//is to write http status code by default, in case it's forgotten in the implementation supplied
	r.httpStatusCodeHandlers.handlers[statusCode] = func(w http.ResponseWriter, r *http.Request, ai *AdditionalInfo) {
		w.WriteHeader(statusCode)
		handler(w, r, ai)
	}
}

//HandleFunc adds a new http request handler for the pattern defined. You also must choose
//which methods this handler will be responding to.
func (r *HttpRouter) HandleFunc(pattern string, methods []string, handler HandlerFunc) {
	r.mx.Lock()
	defer r.mx.Unlock()

	route, err := r.addRoute(pattern)
	if err != nil {
		panic(err)
	}

	route.methods = methods
	if route.methods == nil || len(route.methods) == 0 {
		panic(fmt.Sprintf("method(s) for pattern \"%s\" is not defined!", pattern))
	}

	route.handler = handler
	if route.handler == nil {
		panic(fmt.Sprintf("handler for pattern \"%s\" is not defined!", pattern))
	}
}

//ServerHTTP serves the requests
func (r *HttpRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//Checks whether the request maker is authenticated
	if r.authenticator != nil && !r.authenticator.Authenticate(req) {
		r.httpStatusCodeHandlers.handlers[http.StatusUnauthorized](w, req, nil)
		return
	}

	//Checks whether incoming request is authorized to access the route
	if r.authorizer != nil && !r.authorizer.Authorize(req) {
		r.httpStatusCodeHandlers.handlers[http.StatusForbidden](w, req, nil)
		return
	}

	//Looking for route withing the defined handlers
	route, variables := r.findRoute(req.URL.Path)

	info := r.defaultAdditionalInfo
	if variables != nil {
		info = newAdditionalInfo()
		info.Variables.data = variables
	}

	//404 Not Found
	if route == nil {
		r.httpStatusCodeHandlers.handlers[http.StatusNotFound](w, req, nil)
		return
	}

	//Checks whether the method of the request is allowed for this handler
	allowedMethod := false
	for i := 0; i < len(route.methods); i++ {
		if route.methods[i] != req.Method {
			continue
		}

		allowedMethod = true

		break
	}

	//405 Method Not Allowed
	if !allowedMethod {
		r.httpStatusCodeHandlers.handlers[http.StatusMethodNotAllowed](w, req, nil)
		return
	}

	route.handler(w, req, info)
}

//===========[FUNCTIONALITY]====================================================================================================

//NewRouter crates a new instance of HttpRouter and returns pointer to it
func NewRouter() *HttpRouter {
	return &HttpRouter{
		staticRoutes:           map[string]*route{},
		variableRoutes:         []*route{},
		httpStatusCodeHandlers: newCustomHttpCodeHandlers(),
		defaultAdditionalInfo:  newAdditionalInfo(),
		mx:                     sync.RWMutex{},
	}
}
