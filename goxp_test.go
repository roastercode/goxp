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
