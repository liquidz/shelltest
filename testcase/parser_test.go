package testcase

import (
	"reflect"
	"testing"
)

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
			Comment:  "foo",
		},
	}

	if !reflect.DeepEqual(expected, ts.Tests) {
		t.Errorf("expected testcase: %v, actual testcase %v", expected, ts.Tests)
	}

}

func TestParseWithRegexpMethod(t *testing.T) {
	sample := `
core@foo ~ $ aa
=~ foo
`
	ts, err := Parse(sample)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := TestCases{
		TestCase{
			Command:  "aa",
			Expected: Assertion{RegexpMethod, "foo"}.ToArray(),
		},
	}

	if !reflect.DeepEqual(expected, ts.Tests) {
		t.Errorf("expected testcase: %v, actual testcase %v", expected, ts.Tests)
	}
}

func TestParseErrorWithRegexpMethod(t *testing.T) {
	_, err := Parse(`
core@foo ~ $ aa
=~ foo(
`)
	if err == nil {
		t.Errorf("parsing regexp error should be occured")
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
