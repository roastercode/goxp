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
	// Patch adds a route for a HTTP PATCH request to the specified matching patten.
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
