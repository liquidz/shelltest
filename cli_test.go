package main

import (
	"bytes"
	"fmt"
	. "github.com/liquidz/shelltest/mock"
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

	expected := "no such file"
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("expected %v to contains %v", errStream.String(), expected)
	}
}

func TestRunLint(t *testing.T) {
	cli, _, _ := setup()

	MockReadFile(ReadFileReturn{`
	$ echo foo
	foo
	`, nil})
	defer ResetMock()

	status := cli.Run([]string{"./shelltest", "-l", "foo"})
	if status != ExitCodeOK {
		t.Errorf("expected %v to eq %v", status, ExitCodeOK)
	}

	if len(cli.suite.Tests) != 1 {
		t.Errorf("expected length is 1, but %d", len(cli.suite.Tests))
	}

	expected := TestCase{
		Command:  "echo foo",
		Expected: Assertion{Method: DefaultMethod, Text: "foo"}.ToArray(),
	}

	if !reflect.DeepEqual(expected, cli.suite.Tests[0]) {
		t.Errorf("expected is %v, but %v", expected, cli.suite.Tests[0])
	}
}

func TestRunWithMultipleFiles(t *testing.T) {
	cli, _, _ := setup()

	MockReadFile(
		ReadFileReturn{`
$ echo foo
foo
		`, nil},
		ReadFileReturn{`
$ echo bar
bar
		`, nil},
	)
	defer ResetMock()

	status := cli.Run([]string{"./shelltest", "-l", "foo", "bar"})
	if status != ExitCodeOK {
		t.Errorf("expected %v to eq %v", status, ExitCodeOK)
	}

	if len(cli.suite.Tests) != 2 {
		t.Errorf("expected length is 2, but %d", len(cli.suite.Tests))
	}

	expected := TestCases{
		TestCase{Command: "echo foo", Expected: Assertion{DefaultMethod, "foo"}.ToArray()},
		TestCase{Command: "echo bar", Expected: Assertion{DefaultMethod, "bar"}.ToArray()},
	}

	if !reflect.DeepEqual(expected, cli.suite.Tests) {
		t.Errorf("expected is %v, but %v", expected, cli.suite.Tests)
	}
}
