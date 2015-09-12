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
