package testcase

import (
	"reflect"
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

func TestTestSuiteMerge(t *testing.T) {
	ts1 := TestSuite{
		Tests:  TestCases{TestCase{Command: "foo"}},
		EnvMap: map[string]string{"foo": "bar"},
	}
	ts2 := TestSuite{
		Tests:  TestCases{TestCase{Command: "bar"}},
		EnvMap: map[string]string{"bar": "baz"},
	}
	ts3 := ts1.Merge(ts2)

	if len(ts3.Tests) != 2 {
		t.Errorf("expected length is 2 but %d", len(ts3.Tests))
	}

	expectedTests := TestCases{ts1.Tests[0], ts2.Tests[0]}
	if !reflect.DeepEqual(expectedTests, ts3.Tests) {
		t.Errorf("expected %v but %v", expectedTests, ts3.Tests)
	}

	expectedMap := map[string]string{"foo": "bar", "bar": "baz"}
	if !reflect.DeepEqual(expectedMap, ts3.EnvMap) {
		t.Errorf("expected %v but %v", expectedMap, ts3.EnvMap)
	}
}
