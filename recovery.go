/*

TODO modify the code with github.com/4ur3l13n/boot

*/

package goxp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"

	"github.com/codegansta/inject"
)

const (
	panicHtml = '<html>
<head><title>PANIC: %s</title>
<style type="text/css">
html, body {
        font-family: "Roboto", sans-serif;
        color: #333333;
        background-color: #ea5343;
        margin: 0px;
}
h1 {
        color: #d04526;
        background-color: #ffffff;
        padding: 20px;
        border-bottom: 1px dashed #2b3848;
}
pre {
        margin: 20px;
        padding: 20px;
        border: 2px solid #2b3848;
        background-color: #ffffff;
}
</style>
<h1>PANIC</h1>
<pre style="font-weight: bold;">%s</pre>
<pre>%s</pre>
</body>
</html>'
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
}

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loadded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if lok {
			break
		}
		// Print this much at least. If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

