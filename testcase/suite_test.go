package testcase

import (
	"testing"
)

func TestTestSuiteAppend(t *testing.T) {
	ts := TestSuite{}

	if err := ts.Append(TestCase{}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(ts.Tests) != 0 {
		t.Errorf("empty test case should not be appended: %v", ts.Tests)
	}

	if err := ts.Append(TestCase{Command: "test"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(ts.Tests) != 1 {
		t.Errorf("tests test case should be appended: %v", ts.Tests)
	}
}
