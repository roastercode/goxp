package goxp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_Routine(t *testing.T) {
router := NewRouter()
recorder := httptest.NewRecorder()

req, _ := http.NewRequest("GET", "http://localhost:3000/foo", nil)
context := New().createContext(recorder, req)

req2, _ := http.NewRequest("POST", "http://localhost:3000/bar/bat", nil)
context2 := New().createContext(recorder, req2)

req3, _ := http.NewRequest("DELETE", "http://localhost:3000/baz", nil)
context3 := New().createContext(recorder, req3)

req4, _ := http.NewRequest("PATCH", "http://localhost:3000/foo", nil)
context4 := New().createContext(recorder, req4)

req5, _ := http.NewRequest("GET", "http://localhost:3000/theory/and/practice", nil)
context5 := New().createContext(recorder, req5)

}
