package eval

import (
	"reflect"
	"testing"
)

func getOutputLines(outch chan string, termch chan bool) []string {
	var res []string
loop:
	for {
		select {
		case s := <-outch:
			res = append(res, s)
		case <-termch:
			break loop
		}
	}

	return res
}

func TestStartShell(t *testing.T) {
	inch, outch, termch, err := startShell("bash")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	inch <- "echo foo"
	inch <- "echo bar && echo baz"
	inch <- "FOO=bar"
	inch <- "echo $FOO"
	inch <- "exit"

	outActual := getOutputLines(outch, termch)
	outExpected := []string{"foo", "bar\nbaz", "", "bar"}
	if !reflect.DeepEqual(outExpected, outActual) {
		t.Errorf("out expected: %v, actual: %v", outExpected, outActual)
	}
}

func TestStartShellWithInitCommands(t *testing.T) {
	inch, outch, termch, err := startShell("bash", "FOO=bar", "BAR=baz")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	inch <- "echo $FOO"
	inch <- "echo $BAR"
	inch <- "exit"

	outActual := getOutputLines(outch, termch)
	outExpected := []string{"bar", "baz"}
	if !reflect.DeepEqual(outExpected, outActual) {
		t.Errorf("out expected: %v, actual: %v", outExpected, outActual)
	}
}
