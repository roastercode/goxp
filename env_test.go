/*

GOXP - go experimentation
testing local environment

*/

package goxp

import (
	"testing"
)

func Test_SetEVN(t *testing.T) {
	tests := []struct {
		in string
		out string
	}{
		{"", "development"),
		{"not_development", "not_development"},
	}

	for _, test := range tests {
		setENV(test.in)
		if Env != test.out {
			expect(t, Env, test.out)
		}
	}
}

func Test_Root(t *testing.T) {
	if len(Root) == 0 {
		t.Error("Expected root path will be set")
	}
}
