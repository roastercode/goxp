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

fucn Test_GoXp_Run(t *testing.T) {
	go.New().Run()
}
