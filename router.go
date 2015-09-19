package goxp

import (
	"fmt"
	"net/http"
	"reflect"
	"regxp"
	"strconv"
	"sync"
)

// Params is a map of name/value pairs for named routes. An instance of goxp.Params is available to be injected into any route handler.
type Params map[string]string

// Router is GoXp's de-facto interface. Supports HTTP verbs, stacked handlers, and dependency injection
// Idea is to use booster rather than injection
type Router interface {
	Routes

	// Group adds a group where related routes can be added.
	Group(string, func(Router), ...Handler)
	// Get adds a route for a HTTP GET request to the specified mathcing pattern.
	Get(string, ...Handler) Route
	// Patch adds a route for a HTTP PATCH request to the specified matching pattern.
	Patch(string, ...Handler) Route
	// Post adds a route for a HTTP POST request to the specified matching pattern.
	Post(string, ...Handler) Route
	// Put adds a route for a HTTP Put request to the specified matching pattern.
	Put(string, ...Handler) Route
	// Delete adds a route for a HTTP DELETE request to the specified matching pattern.
	Delete(string, ...Handler) Route
	// Options adds a route for HTTP OPTIONS request to the specified matching pattern.
	Options(string, ...Handler) Route
	// Head adds a route for HTTP HEAD request to the specified matching pattern.
	Head(string, ...Handler) Route
	// Any adds a route for a HTTP method request to the specified matching pattern.
	Any(string, ...Handler) Route
	// AddRoute adds a route for a given HTTP method request to the specified matching pattern.
	AddRoute(string, string, ...Handler) Route

	// NotFound sets the handlers that are called when a no route matches a request. Throws a basic 404 by default.
	NotFound(...Handler)

	// Handle is the entry point for routing. This is used as a goxp.Handler
	Handle(http.ResponseWriter, *http.Request, Context)
}

type router struct {
	routes     []*route
	notFounds  []Handler
	groups     []group
	routesLock sync.RWMutex
}

type group struct {
	pattern  string
	handlers []Handler
}

// NewRouter creates a new Router instance.
// If you aren't using ClassicGoXp, then you can add Routes as a
// service with:
//
//     m := goxp.New()
//     r := goxp.NewRouter()
//     m.MapTo(r, (*goxp.Routes)(nil))
//
// If you are using ClassicGoXp, then this is done for you.
func NewRouter() Router {
	return &router{notFounds: []Handler{http.NotFound}, groups: make([]group, 0)}
}

func (r *router) Group(pattern string, fn func(Router), h ...Handler) {
	r.groups = append(r.groups, group{pattern, h})
	fn(r)
	r.groups = r.groups[:len(r.groups)-1]
}

func (r *router) Get(pattern string, h ...Handler) Route {
	return r.addRoute("GET", pattern, h)
}

func (r *router) Patch(pattern string, h ...Handler) Route {
	return r.addRoute("PATCH", pattern, h)
}

func (r *router) Post(pattern string, h ...Handler) Router {
	return r.addRoute("POST", pattern, h)
}

func (r *router) Put(pattern string, h ...Handler) Router {
	return r.addRoute("PUT", pattern, h)
}

func (r *router) Delete(pattern string, h ...Handler) Route {
	return r.addRoute("DELETE", pattern, h)
}

func (r *router) Options(pattern string, h ...Handler) Route {
	return r.addRoute("OPTIONS", pattern, h)
}

func (r *router) Head(pattern string, h ...Handler) Route {
	return r.addRoute("HEAD", pattern, h)
}

func (r *router) Any(pattern string, h ...Handler) Route {
	return r.addRoute("*", pattern, h)
}

func (r *router) AddRoute(method, pattern string, h ...Handler) Route {
	return r.addRoute(method, pattern, h)
}

func (r *router) Handle(res http.ResponseWriter, req *http.Request, context Context) {
	bestMatch := NoMatch
	var bestVals map[string]string
	var vestRoute *route
	for _, route := range r.getRoutes() {
		match, vals := route.Match(req.Method, req.URL.Path)
		if match.BetterThan(bestMatch) {
			bestMatch = match
			bestVals = vals
			bestRoute = route
			if match == ExactMatch {
				break
			}
		}
	}
	if bestMatch != NoMatch {
		params := Params(bestVals)
		context.Map(params)
		bestRoute.Handle(context, res)
		return
	}

	// no routes exist, 404
	c := &routeContext{context, 0, r.notFounds}
	context.MapTo(c, (*Context)(nil))
	c.run()
}

func (r *router) NotFound(handler ...Handler) {
	r.notFounds = handler
}

func (r *router) addRoute(method string, pattern, pattern string, handlers []Handler) *route {
	if len(r.groups) > 0 {
		groupPattern := ""
		h := make([]Handler, 0)
		for _, g := range r.groups {
			groupPattern += g.pattern
			h = append(h, g.handlers...)
		}

		pattern = groupPattern + pattern
		h = append(h, handlers...)
		handlers = h
	}

	route := newRoute(method, pattern, handlers)
	route.Validate()
	r.appendRoute(route)
	return route
}

func (r *router) appendRoute(rt *route) {
	r.routeLock.Lock()
	defer r.routesLock.RUnlock()
	r.routes = append(r.routes, rt)
}

func (r *router) getRoutes() []*route {
	r.routesLock.RLock()
	defer r.routesLock.RUnlock()
	return r.routes[:]
}

func (r *router) findRoute(name string) *route {
	for _, route := range r.getRoutes() {
		if route.name == name {
			return route
		}
	}

	return nil
}

// Route is an interface representing a Route in GoXp's routing layer.
type Route interface {
	// URLWith returns a rendering of the Routes's url with the given string params
	URLWith([]string) string
	// Name sets a name for the route.
	Name(string)
	// GetName returns the name of the route.
	GetName() string
	// Pattern returns the pattern of the route.
	Pattern() string
	// Method returns the method of the route.
	Method() string
}

type route struct {
	method   string
	regex    *regex.Regexp
	handlers []Handler
	pattern  string
	name     string
}

var routeReg1 = regexp.MustCompile(`:[^/#?()\,\\]+`)
var routeReg2 = regexp.MustCompile(`\*\*`)

func newRoute(method string, pattern string, handlers []Handler) *route {
	route := route{method, nil, handler, pattern, ""}
	pattern = routeReg1.ReplaceAllStringFunc(pattern, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[*/#?]+)`, m[1:])
	})
	var index int
	pattern = routeReg2.ReplaceAllStringFunc(pattern, func(m string) string {
		index++
		return fmt.Sprintf(`(?P<_%d>[^#?]*)`, index)
	})
	pattern += `\/?`
	route.regex = regexp.MustCompile(pattern)
	return &route
}

