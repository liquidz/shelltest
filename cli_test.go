package main

import (
	"bytes"
	"fmt"
	. "github.com/liquidz/shelltest/testcase"
	"reflect"
	"strings"
	"testing"
)

func setup() (*CLI, *bytes.Buffer, *bytes.Buffer) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}

	return cli, outStream, errStream
}

func TestRunVersionFlag(t *testing.T) {
	cli, _, errStream := setup()

	status := cli.Run([]string{"./shelltest", "-v"})
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	expected := fmt.Sprintf("shelltest version %s", Version)
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("expected %q to eq %q", errStream.String(), expected)
	}
}

func TestRunNoArgumentError(t *testing.T) {
	cli, _, errStream := setup()

	status := cli.Run([]string{"./shelltest"})
	if status != ExitCodeError {
		t.Errorf("expected %v to eq %v", status, ExitCodeError)
	}

	expected := "no arg"
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("expected %v to contains %v", errStream.String(), expected)
	}
}

func TestRunNotExistingFile(t *testing.T) {
	cli, _, errStream := setup()

	status := cli.Run([]string{"./shelltest", "not_existing_file"})
	if status != ExitCodeError {
		t.Errorf("expected %v to eq %v", status, ExitCodeError)
	}

	expected := "not found"
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("expected %v to contains %v", errStream.String(), expected)
	}
}

func TestRunLint(t *testing.T) {
	cli, _, _ := setup()

	MockReadFile(`
	$ echo foo
	foo
	`, nil)
	defer ResetMock()

	status := cli.Run([]string{"./shelltest", "-l", "foo"})
	if status != ExitCodeOK {
		t.Errorf("expected %v to eq %v", status, ExitCodeOK)
	}

	if len(cli.suite.Tests) != 1 {
		t.Errorf("FIXME")
	}

	expected := TestCase{
		Command:  "echo foo",
		Expected: []Assertion{Assertion{Method: DefaultMethod, Text: "foo"}},
	}

	if !reflect.DeepEqual(expected, cli.suite.Tests[0]) {
		t.Errorf("FIXME")
	}
}

//func TestRun_fFlag(t *testing.T) {
//	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
//	cli := &CLI{outStream: outStream, errStream: errStream}
//	args := strings.Split("./shunig -f", " ")
//
//	status := cli.Run(args)
//	_ = status
//}
