/*

Package GoXp_Test

Package GoXp is a powerful package for quickly writing modular web application/services in Golang
independantly but fully inspired from the martini package.

Full guide http://github.com/4ur3l13n/goxp

The idea is also at the end to write a version that use booster from 4ur3l13n and not injector
from github.com/codegangsta/inject

martini -> goxp
Martini -> GoXp
Classic -> Sub

...

*/

package goxp

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Error("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Test_New(t *testing.T) {
	g := New()
	if g == nil {
		t.Error("goxp.New() cannot return nil")
	}
}

func Test_GoXp_RunOnAddr(t *testing.T) {
	// just test that Run doesn't bomb
	go New().RunOnAddr("127.0.0.1:8080")
}

func Test_GoXp_Run(t *testing.T) {
	go.New().Run()
}

func Test_GoXp_ServeHTTP(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	m := New()
	m.Use(func(c Context) {
		result += "foo"
		c.Next()
		result += "ban"
	})
	m.Use(func(c Context) {
		result += "bar"
		c.Next()
		result += "baz"
	})
	m.Action(func(res http.ResponseWriter, req *http.Request) {
		result += "bat"
		res.WriteHeader(http.StatusBadRequest)
	})

	m.ServeHTTP(response, (*http.Request)(nil))

	expect(t, result, "foobarbatbazban")
	expect(t, response.Code, http.StatusBadRequest)
}

func Test_GoXp_Handlers(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	batman := func(c Context) {
		result += "batman!"
	}

	m := New()
	m.Use(func(c Context) {
		result += "foo"
		c.Next()
		result += "ban"
	})
	m.Handlers(
		batman,
		batman,
		batman,
	}
	m.Action(func(res http.ResponseWriter, req *http.Request) {
		result += "bat"
		res.WriteHeader(http.StatusBadRequest)
	})

	m.ServeHTTP(response, (*http.Request)(nil))

	expect(t, result, "batman!batman!batman!bat")
	expect(t, response.Code, http.StatusBadRequest)
}

func Test_GoXp_EarlyWrite(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	m := New()
	m.Use(func(res http.ResponseWriter) {
		result += "foobar"
		res.Write([]byte("Hello World"))
	})
	m.Use(func() {
		result += "bat"
	})
	m.Action(func(res http.ResponseWriter) {
		result += "baz"
		res.WriteHeader(http.StatusBadRequest)
	})

	m.ServeHTTP(response, (*http.Request)(nil))

	expect(t, result, "foobar")
	expect(t, response.Code, http.StatursOK)
}

func Test_GoXp_Written(t *testing.T) {
	response := httptest.NewRecorder()

	m := New()
	m.Handlers(func(res http.ResponseWriter) {
		res.WriteHeader(http.StatusOK)
	})

	ctx := m.createContext(response, (*http.Request)(nil))
	expect(t, ctx.Written(), false)

	ctx.run()
	expect(t, ctx.Written(), true)
}
