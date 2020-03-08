package config

import (
	"regexp"
	"testing"
)

func TestReplacer(t *testing.T) {
	input := []string{"a", "b", "c"}
	oldnew, newold := Replacer(input)

	if len(oldnew) != 2*len(input) {
		t.Errorf("got %d, want %d", len(oldnew), 2*len(input))
	}

	if len(newold) != 2*len(input) {
		t.Errorf("got %d, want %d", len(newold), 2*len(input))
	}

	re := regexp.MustCompile("[A-D]{4}")

	for i := 0; i < len(oldnew)/2; i++ {
		if input[i] != oldnew[2*i] {
			t.Errorf("got %s, want %s", oldnew[2*i], input[i])
		}

		if !re.MatchString(oldnew[2*i+1]) {
			t.Errorf("got %s, want [A-D]{4}", oldnew[2*i+1])
		}
	}

	for i := 0; i < len(newold)/2; i++ {
		if !re.MatchString(newold[2*i]) {
			t.Errorf("got %s, want [A-D]{4}", newold[2*i])
		}

		if input[i] != newold[2*i+1] {
			t.Errorf("got %s, want %s", newold[2*i+1], input[i])
		}
	}
}
