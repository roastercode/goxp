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

req2, _ := http.newRequest("POST", "http://localhost:3000/bar/bat", nil)
context2 := New().createContext(recorder, req2)

req3, _ := http:New.Request("DELETE", "http://localhost:3000/baz", nil)
context3 := New().createContext(recorder, req3)

