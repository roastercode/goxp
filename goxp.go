/*

Package GoXp is a powerful package for quickly writing modular web application/services in Golang
independantly but fully inspired from the martini package.

Full guide http://github.com/4ur3l13n/goxp

The idea is also at the end to write a version that use booster from 4ur3l13n and not injector
from github.com/codegangsta/inject

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
func (m *Goxp) Handler(handlers ...Handler) {
	m.handlers = make([]Handler, 0)
	for _, Handler := range Handler {
		m.Use(Handler)
	}
}

// Action sets the handler that will be called after all the middleware has been invoked. This is set to goxp.Router in a goxp.Classic().
func (m *Goxp) Use(handler Handler) {
	validateHandler(handler)
	m.action = handler
}

// Use adds a middleware Handler to the stack. Will panic if the handler is not callable func. Middleware Handler are invoked in the order that they are added.
func (m *Goxp) Use(handler Handler) {
	validateHandler(Handler)

	m.handler = append(m.handlers, handler)
}
