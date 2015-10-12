package main

import (
	"reflect"
	"testing"
)

func TestAssertionAssert(t *testing.T) {
	var a Assertion

	a = Assertion{Method: EqualMethod, Text: "foo"}
	tests := map[string]bool{
		"foo": true,
		"bar": false,
	}
	for s, expected := range tests {
		if a.Assert(s) != expected {
			t.Errorf("assert %s should be %v", s, expected)
		}
	}
}

func TestAssertionsIsExpected(t *testing.T) {
	var as Assertions
	tests := map[string]bool{
		"":              false,
		"foo":           false,
		"foo\nbar":      true,
		"foo\nbar\n":    true,
		"foo\nbar\nbaz": true,
	}

	as = Assertions{
		Assertion{EqualMethod, "foo"},
		Assertion{EqualMethod, "bar"},
	}
	for s, expected := range tests {
		if as.IsExpected(s) != expected {
			t.Errorf("assert [%s] should be %v", s, expected)
		}
	}

	as = Assertions{}
	for s, _ := range tests {
		if !as.IsExpected(s) {
			t.Errorf("assert %s should be true", s)
		}
	}
}

func TestTestCaseIsEmpty(t *testing.T) {
	var tc TestCase

	tc = TestCase{}
	if !tc.IsEmpty() {
		t.Errorf("test case should be empty")
	}

	tc = TestCase{Command: "foo"}
	if tc.IsEmpty() {
		t.Errorf("test case should not be empty: %v", tc)
	}

	as := Assertions{Assertion{"foo", "bar"}}

	tc = TestCase{Expected: as}
	if tc.IsEmpty() {
		t.Errorf("test case should not be empty: %v", tc)
	}

	tc = TestCase{Command: "foo", Expected: as}
	if tc.IsEmpty() {
		t.Errorf("test case should not be empty: %v", tc)
	}
}

func TestTestCaseAppendAssertion(t *testing.T) {
	var tc TestCase

	tc = TestCase{}
	tc.AppendAssertion("foo", "bar")
	if len(tc.Expected) != 0 {
		t.Errorf("assertion should not be appended to empty test case: %v", tc.Expected)
	}

	tc = TestCase{Command: "foo"}
	tc.AppendAssertion("foo  ", "bar")
	tc.AppendAssertion("bar", "  baz   ")

	expected := TestCase{
		Command: "foo",
		Expected: Assertions{
			Assertion{"foo", "bar"},
			Assertion{"bar", "baz"},
		},
	}

	if !reflect.DeepEqual(expected, tc) {
		t.Errorf("expected test case: %v, actual test case: %v", expected, tc)
	}
}

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

func TestParseWithDefaultSection(t *testing.T) {
	sample := `
core@foo ~ $ command
foo
	`

	ts, err := Parse(sample)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(ts.Tests) != 1 {
		t.Errorf("test: expected length: 1, actual length: %v", len(ts.Tests))
	}

	expected := TestCase{
		Command:  "command",
		Expected: Assertion{DefaultMethod, "foo"}.ToArray(),
	}

	if !reflect.DeepEqual(expected, ts.Tests[0]) {
		t.Errorf("expected testcase: %v, actual testcase %v", expected, ts.Tests[0])
	}
}

func TestParseWithMultipleLineCommand(t *testing.T) {
	sample := `
core@foo ~ $ aa \
bb
cc
`
	ts, err := Parse(sample)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := TestCases{
		TestCase{
			Command:  "aa bb",
			Expected: Assertion{DefaultMethod, "cc"}.ToArray(),
		},
	}

	if !reflect.DeepEqual(expected, ts.Tests) {
		t.Errorf("expected testcase: %v, actual testcase %v", expected, ts.Tests)
	}
}

func TestParseWithCommentLine(t *testing.T) {
	sample := `
# foo
core@foo ~ $ aa
  # bar
bb
	# baz
`
	ts, err := Parse(sample)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := TestCases{
		TestCase{
			Command:  "aa",
			Expected: Assertion{DefaultMethod, "bb"}.ToArray(),
		},
	}

	if !reflect.DeepEqual(expected, ts.Tests) {
		t.Errorf("expected testcase: %v, actual testcase %v", expected, ts.Tests)
	}

}

func TestParseWithSpecifiedSection(t *testing.T) {
	sample := `
[before]
core@foo ~ $ command

[after]
core@foo ~ $ after1
core@foo ~ $ after2

[test]
core@foo ~ $ echo foo
foo
core@foo ~ $ ls
a  b
cc ddd
	`

	ts, err := Parse(sample)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(ts.Before) != 1 {
		t.Errorf("before: expected length = 1, actual length = %v", len(ts.Before))
	}

	if len(ts.After) != 2 {
		t.Errorf("after: expected length = 2, actual length = %v", len(ts.After))
	}

	if len(ts.Tests) != 2 {
		t.Errorf("test: expected length: 2, actual length: %v", len(ts.Tests))
	}

	expected := TestCases{
		TestCase{
			Command:  "echo foo",
			Expected: Assertion{DefaultMethod, "foo"}.ToArray(),
		},
		TestCase{
			Command: "ls",
			Expected: Assertions{
				Assertion{DefaultMethod, "a  b"},
				Assertion{DefaultMethod, "cc ddd"},
			},
		},
	}

	if !reflect.DeepEqual(expected, ts.Tests) {
		t.Errorf("expected testcase: %v, actual testcase %v", expected, ts.Tests[0])
	}
}
