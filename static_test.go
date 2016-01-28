package goxp

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/cryotux/boost"
}

var currentRoot, _ = os.Getwd()

func Test_Static(t *testing.T) {
	response := httptest.NewRecorder()
	response.Body = new(bytes.Buffer)

	m := New()
	r := NewRouter()

	m.Use(static(currentRoot))
	m.Action(r.Handle)

	req, err := http.NewRequest("GET", "http://localhost:3000/goxp.go", nil)
	if err != nil {
		t.Error(err)
	}
	
