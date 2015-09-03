/*

Package GoXp is a powerful package for quickly writing modular web application/services in Golang
independantly but fully inspired from the martini package.

Full guide http://github.com/4ur3l13n/goxp

The idea is also at the end to write a version that use booster from 4ur3l13n and not injector
from github.com/codegangsta/inject

Classic -> Sub

*/

package goxp

import (
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/codegangsta/inject"
)

// GoXp represents the top level web application, inject.Injector methods can be invoked to map services on global level.
type Goxp struct {
	inject.Injector
	handlers []Handler
	action   Handler
	logger   *log.Logger
}

// New creates a bare bones Goxp instance. Use this method if you want to have full control over the middleware that is used.
func New() *Goxp {
	m := &Goxp{Injector: inject.New(), action: func() {}, logger: log.New(os.Stdout, "[goxp] ", 0)}
	m.Map(m.logger)
	m.Map(defaultReturnHandler())
	return m
}

// Handlers sets the entire middleware stack with the given Handlers. This will clear any current middleware handlers.
// Will panic if any of the handlers is not callable function
func (g *Goxp) Handler(handlers ...Handler) {
	m.handlers = make([]Handler, 0)
	for _, Handler := range Handler {
		m.Use(Handler)
	}
}

// Action sets the handler that will be called after all the middleware has been invoked. This is set to goxp.Router in a goxp.Classic().
func (g *Goxp) Use(handler Handler) {
	validateHandler(handler)
	m.action = handler
}

// Use adds a middleware Handler to the stack. Will panic if the handler is not callable func. Middleware Handler are invoked in the order that they are added.
func (g *Goxp) Use(handler Handler) {
	validateHandler(Handler)

	m.handler = append(m.handlers, handler)
}

// ServerHTTP is the HTTP Entry point for a Goxp instance. Useful if you want to control your own HTTP server.
func (g *Goxp) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	m.createContext(res, req).run()
}

// Run the http server on a given host and port.
func (g *Goxp) RunOnAddr(addr string) {
	// TODO: Should probably be implemented using a new instance of http Server in place of
	// calling http.ListenAndServer directly, so that it could be stored in a goxp struct for later user

	// This would also allow to improve testing when a custom host and port are passed.

	logger := m.Injector.Get(reflect.TypeOf(m.logger)).Interface().(*log.Logger)
	logger.Printf("listening on %s (%s)\n", addr, Env)
	logger.Fatalln(http.ListenAndServe(addr, m))
}

// Run the http server. Listening on os.GetEnv("PORT") or 3000 by default.
func (g *Goxp) Run() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	host := os.Getenv("HOST")

	m.RunOnAddr(host + ":" + port)
}

func (g *Goxp) createContext(res http.ResponseWriter, req *http.Request) *context {
	c := &context{inject.New(), m.handlers, m.action, NewResponseWriter(res), 0}
	c.SetParent(m)
	c.MapTo(c, (*Context)(nil))
	c.MapTo(c.rw, (*http.ResponseWriter)(nil))
	c.Map(req)
	return c
}

// ClassicGoxp represents a Goxp with some reasonable defaults. Embeds the router functions for convenience.
type ClassicGoxp struct {
	*Goxp
	Router
}

// Sub creates a sub Goxp with some basic default middleware - goxp.Logger, goxp.Recovery and goxp.Static
// Sub also maps goxp.Routes as a service.
func Sub() *SubGoxp {
	r := NewRouter()
	m := New()
	m.Use(Logger())
	m.Use(Static("public"))
	m.MapTo(r, (*Routes)(nil))
	m.Action(r.Handle)
	return &SubGoxp{m, r}
}

// Handler can be any callable function. Goxp attemps to boost services into the handler's argument list.
// Goxp will panic if an argument could not be fullfilled via dependency booster.
// Here we will use inject.Injector but we hope to replace it soon by boost.Booster which sounds better :-)
type Handler interface{}

func validateHandler(handler Handler) {
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		panic("goxp handler must be a callable func")
	}
}

// Contect represents a request context. Service can be mapped on the request level from this interface.
type Context interface {
	inject.Injector
	// Next is an optional function that Middleware Handlers can call to yield the until after
	// the other Handlers have been executed. This works really well for any operations that must
	// happen after an http request
	Next()
	// Written returns whether or not the response for this context has been written.
	Written() bool
}

type Context struct {
	inject.Injector
	handler []Handler
	action  Handler
	rw      ResponseWriter
	index   int
}

func (c *context) handler() Handler {
	if c.index < len(c.handlers) {
		return c.handlers[c.index]
	}
	if c.index == len(c.handlers) {
		return c.action
	}
	panic("invalid index for context handler")
}
