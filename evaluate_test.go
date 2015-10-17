package main

import (
	"testing"
)

func nocallback(_ TestCase, _ error) {}

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

	errs := Evaluate(DefaultShell, ts, nocallback)
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
	errs := Evaluate(DefaultShell, ts, nocallback)
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
	errs := Evaluate(DefaultShell, ts, nocallback)
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
	errs := Evaluate(DefaultShell, ts, nocallback)
	if len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
}

func TestEvaluateContainingEmptyOutput(t *testing.T) {
	ts := TestSuite{
		Tests: TestCases{
			TestCase{"FOO=bar", Assertions{}},
			TestCase{"echo $FOO", Assertions{Assertion{DefaultMethod, "foo"}}},
		},
	}
	errs := Evaluate(DefaultShell, ts, nocallback)
	if len(errs) != 1 {
		t.Errorf("one error should be occured: %v", errs)
	}
}
