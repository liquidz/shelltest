package main

import (
	"reflect"
	"testing"
)

func TestStartShell(t *testing.T) {
	inch, outch, termch, err := startShell("bash")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	inch <- "echo foo"
	inch <- "echo bar && echo baz"
	inch <- "exit"

	var outActual []string
loop:
	for {
		select {
		case s := <-outch:
			outActual = append(outActual, s)
		case <-termch:
			break loop
		}
	}

	outExpected := []string{"foo", "bar\nbaz"}
	if !reflect.DeepEqual(outExpected, outActual) {
		t.Errorf("out expected: %v, actual: %v", outExpected, outActual)
	}
}
