/*

GOXP - go experimentation

*/

package goxp

import (
	"os"
)

// Envirronment
const (
	Dev  string = "development"
	Prod string = "production"
	Test string = "test"
)

// env is the environment that goxp is executing in. The GOXP_ENV is read on initialization
// to set variable

var Env = Dev
var Root string

func setENV(e string) {
	if len(e) > 0 {
		Env = e
	}
}

func init() {
	setENV(os.Getenv("GOXP_ENV"))
	var err error
	Root, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}
