package formatter

import (
	. "github.com/liquidz/shelltest/eval"
	. "github.com/liquidz/shelltest/testcase"
	"strings"
	"testing"
)

func TestDefaultSetup(t *testing.T) {
	f := DefaultFormatter{}
	actual := f.Setup(TestSuite{})
	if actual != "" {
		t.Errorf("response should be empty string: %s", actual)
	}
}

func TestDefaultResultSuccess(t *testing.T) {
	f := DefaultFormatter{}
	var expected, actual string

	actual = f.Result(0, TestCase{}, nil)
	expected = "."
	if expected != actual {
		t.Errorf("response should be %v but %v", expected, actual)
	}
}

func TestDefaultResultFail(t *testing.T) {
	f := DefaultFormatter{}
	tc := TestCase{Command: "foo"}
	var expected, actual string

	actual = f.Result(0, tc, EvaluateError{No: 1, Test: tc})
	expected = "1) foo"
	if !strings.Contains(actual, expected) {
		t.Errorf("response should contains %v but %v", expected, actual)
	}

	tc = TestCase{Command: "foo", Comment: "bar"}
	actual = f.Result(0, tc, EvaluateError{No: 2, Test: tc})
	expected = "2) bar"
	if !strings.Contains(actual, expected) {
		t.Errorf("response should contains %v but %v", expected, actual)
	}
}

func TestDefaultTearDown(t *testing.T) {
	f := DefaultFormatter{}
	tc := TestCase{Expected: Assertion{"a", "b"}.ToArray()}
	ts := TestCases{tc, tc}
	var expected, actual string

	actual = f.TearDown(TestSuite{Tests: ts}, []error{})
	expected = "2 tests, 0 failures"
	if !strings.Contains(actual, expected) {
		t.Errorf("response should contains %v but %v", expected, actual)
	}

	actual = f.TearDown(TestSuite{Tests: ts}, []error{EvaluateError{}})
	expected = "2 tests, 1 failures"
	if !strings.Contains(actual, expected) {
		t.Errorf("response should contains %v but %v", expected, actual)
	}
}
