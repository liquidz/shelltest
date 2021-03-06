package eval

import (
	. "github.com/liquidz/shelltest/testcase"
	"testing"
)

func nocallback(_ int, _ TestCase, _ error) {}

func TestEvaluate(t *testing.T) {
	ts := TestSuite{
		Tests: TestCases{
			TestCase{
				Command: "echo foo && echo bar",
				Expected: []Assertion{
					Assertion{DefaultMethod, "foo"},
					Assertion{DefaultMethod, "bar"},
				},
			},
		},
	}

	errs := Evaluate("bash", ts, nocallback)
	if len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
}

func TestEvaluateInvalidOutputError(t *testing.T) {
	ts := TestSuite{
		Tests: TestCases{
			TestCase{
				Command: "echo foo",
				Expected: Assertions{
					Assertion{DefaultMethod, "bar"},
				},
			},
		},
	}
	errs := Evaluate("bash", ts, nocallback)
	if len(errs) != 1 {
		t.Errorf("unexpected errors: %v", errs)
	}
}

func TestEvaluateTooFewOutputError(t *testing.T) {
	ts := TestSuite{
		Tests: TestCases{
			TestCase{
				Command: "echo foo",
				Expected: Assertions{
					Assertion{DefaultMethod, "foo"},
					Assertion{DefaultMethod, "bar"},
				},
			},
		},
	}
	errs := Evaluate("bash", ts, nocallback)
	if len(errs) != 1 {
		t.Errorf("unexpected errors: %v", errs)
	}
}

func TestEvaluateWithReturnCode(t *testing.T) {
	ts := TestSuite{
		Tests: TestCases{
			TestCase{
				Command:  "(exit 10)",
				Expected: []Assertion{},
			},
			TestCase{
				Command:  "echo $?",
				Expected: []Assertion{Assertion{DefaultMethod, "10"}},
			},
		},
	}
	errs := Evaluate("bash", ts, nocallback)
	if len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
}

func TestEvaluateContainingEmptyOutput(t *testing.T) {
	ts := TestSuite{
		Tests: TestCases{
			TestCase{"FOO=bar", Assertions{}, ""},
			TestCase{"echo foo$FOO", Assertions{Assertion{DefaultMethod, "baz"}}, ""},
		},
	}
	errs := Evaluate("bash", ts, nocallback)
	if len(errs) != 1 {
		t.Errorf("one error should be occured: %v", errs)
	}
}

func TestEvaluateWithEnvMap(t *testing.T) {
	ts := TestSuite{
		Tests: TestCases{
			TestCase{"echo foo$FOO", Assertions{Assertion{DefaultMethod, "foobar"}}, ""},
		},
		EnvMap: map[string]string{
			"FOO": "bar",
		},
	}

	errs := Evaluate("bash", ts, nocallback)
	if len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
}

// TODO: Evaluate cannot receive empty string from outch
//func TestEvaluateWithUndefinedVariable(t *testing.T) {
//	ts := TestSuite{
//		Tests: TestCases{
//			TestCase{"echo $FOO", Assertions{Assertion{DefaultMethod, "foobar"}}, ""},
//		},
//	}
//
//	errs := Evaluate("bash", ts, nocallback)
//	if len(errs) != 0 {
//		t.Errorf("unexpected errors: %v", errs)
//	}
//}
