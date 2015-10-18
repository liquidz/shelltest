package formatter

import (
	. "github.com/liquidz/shelltest/color"
	. "github.com/liquidz/shelltest/eval"
	. "github.com/liquidz/shelltest/testcase"
	"regexp"
	"strings"
	"testing"
)

func TestTapSetup(t *testing.T) {
	f := TapFormatter{}
	var expected, actual string

	actual = f.Setup(TestSuite{})
	expected = "0..0"
	if !strings.Contains(actual, expected) {
		t.Errorf("response should contains %v but %v", expected, actual)
	}

	actual = f.Setup(TestSuite{Tests: TestCases{TestCase{}, TestCase{}}})
	expected = "1..2"
	if !strings.Contains(actual, expected) {
		t.Errorf("response should contains %v but %v", expected, actual)
	}
}

func TestTapResultSuccess(t *testing.T) {
	NoColor = true
	f := TapFormatter{}

	actual := f.Result(0, TestCase{Command: "foo"}, nil)
	expected := regexp.MustCompile(`^ok 1 - foo`)
	if !expected.MatchString(actual) {
		t.Errorf("response should match %v but %v", expected, actual)
	}
}

func TestTapResultFail(t *testing.T) {
	NoColor = true
	f := TapFormatter{}
	tc := TestCase{Command: "foo", Comment: "bar"}

	actual := f.Result(1, tc, EvaluateError{No: 1, Test: tc})
	expected := regexp.MustCompile(`^not ok 2 - bar`)
	if !expected.MatchString(actual) {
		t.Errorf("response should match %v but %v", expected, actual)
	}
}

func TestTapTearDown(t *testing.T) {
	f := TapFormatter{}
	actual := f.TearDown(TestSuite{}, []error{})
	if actual != "" {
		t.Errorf("response should be empty but %v", actual)
	}
}
