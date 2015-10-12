package main

import (
	"testing"
)

func TestStartShell(t *testing.T) {
	in, out, term, err := startShell("bash")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	in <- "echo foo"
	in <- "echo bar && echo baz"
	in <- "exit"

	expected := []string{
		"foo",
		"bar\nbaz",
	}
loop:
	for i := 0; ; i++ {
		select {
		case s := <-out:
			if s != expected[i] {
				t.Errorf("expected: foo, actual: [%v]", expected[i])
			}
		case <-term:
			break loop
		}
	}
}
