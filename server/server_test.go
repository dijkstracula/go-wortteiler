package server

import (
	"strings"
	"testing"
)

func TestValidateInput(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		// valid
		{"abc", true},
		{"klärung", true},
		{"Ärger", true},
		{"Österreich", true},
		{"Übel", true},
		{"Maßen", true},

		// invalid
		{"", false},
		{"number1", false},
		{"<script>", false},
		{strings.Repeat("a", 65), false},
	}

	for _, test := range tests {
		if actual := validateInput(test.in); actual != test.out {
			t.Errorf("validateInput(%+q): got %v, expected %v", test.in, actual, test.out)
		}
	}
}
