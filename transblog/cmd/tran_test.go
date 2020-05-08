package cmd

import (
	"strings"
	"testing"
)

func Test_transTitleLine(t *testing.T) {
	tests := []struct {
		source []byte
		dest   []byte
		wanted bool
	}{
		{[]byte("**** This is a title"), []byte("#### This is a title"), true},
		{[]byte("**** This is a title"), []byte("### This is a title"), false},
	}

	for _, tt := range tests {
		mdTitle := transTitleLine(tt.source)
		cmpRes := strings.Compare(string(tt.dest), string(mdTitle))
		if cmpRes == 0 && tt.wanted {
			t.Log("test ok!")
		} else if cmpRes != 0 && !tt.wanted {
			t.Log("test ok!")
		} else {
			t.Error("test error!")
		}
	}
}

func Test_isOrgSource(t *testing.T) {
	tests := []struct {
		source []byte
		wanted bool
	}{
		{[]byte("#+BEGIN_SRC go"), true},
		{[]byte("#+END_SRC go"), true},
		{[]byte("#+END_SRT go"), false},
	}

	for _, tt := range tests {
		isSource := isOrgSouce(tt.source)
		if isSource != tt.wanted {
			t.Errorf("line %s is not a source line, but wanted %v\n", string(tt.source), tt.wanted)
		}
	}
}

func Test_transSourceLine(t *testing.T) {
	tests := []struct {
		source []byte
		dest   []byte
		wanted bool
	}{
		{[]byte("#+BEGIN_SRC go"), []byte("``` go"), true},
		{[]byte("#+END_SRC"), []byte("```"), true},
	}

	for _, tt := range tests {
		mdSrcLine := transSourceLine(tt.source)
		cmpRes := strings.Compare(string(mdSrcLine), string(tt.dest))
		if cmpRes == 0 && tt.wanted {
			t.Log("test ok!")
		}

		if cmpRes != 0 && !tt.wanted {
			t.Error("test error!")
		}
	}
}
