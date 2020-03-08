package cli

import (
	"testing"
)

func TestTrimUnnecessary(t *testing.T) {
	trimtests := []struct {
		in  string
		out string
	}{
		{`// I am
		// a Gopher.`, "I am a Gopher."},
		{`  # aaa
	# bbb`, "aaa bbb"},
		{`I
am a Gopher.

I like Go.`, "I am a Gopher. \n I like Go."},
	}

	for _, tt := range trimtests {
		t.Run(tt.in, func(t *testing.T) {
			s := trimUnnecessary(tt.in)
			if s != tt.out {
				t.Errorf("got %q, want %q", s, tt.out)
			}
		})
	}
}

func TestAutoDetectJP(t *testing.T) {
	trimtests := []struct {
		in  string
		out bool
	}{
		{"私はGopherです", true},
		{"I am a Gopher.", false},
		{"I am a 社員", false},
	}

	for _, tt := range trimtests {
		t.Run(tt.in, func(t *testing.T) {
			bl := autoDetectJP(tt.in)
			if bl != tt.out {
				t.Errorf("got %t, want %t", bl, tt.out)
			}
		})
	}
}
