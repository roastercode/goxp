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
