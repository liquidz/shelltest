package testcase

import (
	"testing"
)

func TestTestSuiteAppendError(t *testing.T) {
	ts := TestSuite{}

	if err := ts.Append("unknown section", TestCase{Command: "foo"}); err == nil {
		t.Errorf("error should be occured")
	}
}

func TestTestSuiteAppend(t *testing.T) {
	ts := TestSuite{}

	if err := ts.Append("test", TestCase{}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(ts.Tests) != 0 {
		t.Errorf("empty test case should not be appended: %v", ts.Tests)
	}

	if err := ts.Append("before", TestCase{Command: "before"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(ts.Before) != 1 {
		t.Errorf("before test case should be appended: %v", ts.Before)
	}

	if err := ts.Append("after", TestCase{Command: "after"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(ts.After) != 1 {
		t.Errorf("after test case should be appended: %v", ts.After)
	}

	if err := ts.Append("test", TestCase{Command: "test"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(ts.Tests) != 1 {
		t.Errorf("tests test case should be appended: %v", ts.Tests)
	}
}
