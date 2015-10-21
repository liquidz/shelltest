package testcase

import (
	. "github.com/liquidz/shelltest/mock"
	"reflect"
	"testing"
)

func TestParseWithDefaultSection(t *testing.T) {
	MockReadFile(ReadFileReturn{`
core@foo ~ $ command
foo
	`, nil})
	defer ResetMock()

	ts, err := Parse("foo")
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
	MockReadFile(ReadFileReturn{`
core@foo ~ $ aa \
bb
cc
	`, nil})
	defer ResetMock()

	ts, err := Parse("foo")
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
	MockReadFile(ReadFileReturn{`
# foo
core@foo ~ $ aa
  # bar
bb
	# baz
	`, nil})
	defer ResetMock()

	ts, err := Parse("foo")
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
	MockReadFile(ReadFileReturn{`
core@foo ~ $ aa
=~ foo
	`, nil})
	defer ResetMock()

	ts, err := Parse("foo")
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
	MockReadFile(ReadFileReturn{`
core@foo ~ $ aa
=~ foo(
	`, nil})
	defer ResetMock()

	_, err := Parse("foo")
	if err == nil {
		t.Errorf("parsing regexp error should be occured")
	}
}

func TestParseWithAutoAssertion(t *testing.T) {
	MockReadFile(ReadFileReturn{`
core@foo ~ $ aa
core@foo ~ $ bb
	`, nil})
	defer ResetMock()

	ts, err := Parse("foo")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := TestCases{
		TestCase{Command: "aa"},
		TestCase{Command: "echo $?", Expected: Assertion{DefaultMethod, "0"}.ToArray(), Comment: "aa"},
		TestCase{Command: "bb"},
		TestCase{Command: "echo $?", Expected: Assertion{DefaultMethod, "0"}.ToArray(), Comment: "bb"},
	}
	if !reflect.DeepEqual(expected, ts.Tests) {
		t.Errorf("expected testcase: %v, actual testcase %v", expected, ts.Tests)
	}
}

func TestParseWithoutAutoAssertion(t *testing.T) {
	MockReadFile(ReadFileReturn{`
core@foo ~ $ aa
core@foo ~ $ bb
	`, nil})
	defer ResetMock()

	NoAutoAssertion = true
	ts, err := Parse("foo")
	NoAutoAssertion = false
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := TestCases{
		TestCase{Command: "aa"},
		TestCase{Command: "bb"},
	}

	if !reflect.DeepEqual(expected, ts.Tests) {
		t.Errorf("expected testcase: %v, actual testcase %v", expected, ts.Tests)
	}
}

func TestParseWithRequire(t *testing.T) {
	MockReadFile(
		ReadFileReturn{`
$ foo
bar
@require outerfile.txt
$ bar
baz
		`, nil},
		ReadFileReturn{`
$ hello
world
		`, nil},
	)
	defer ResetMock()

	ts, err := Parse("foo")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := TestCases{
		TestCase{Command: "foo", Expected: Assertion{DefaultMethod, "bar"}.ToArray()},
		TestCase{Command: "hello", Expected: Assertion{DefaultMethod, "world"}.ToArray()},
		TestCase{Command: "bar", Expected: Assertion{DefaultMethod, "baz"}.ToArray()},
	}

	if !reflect.DeepEqual(expected, ts.Tests) {
		t.Errorf("expected testcase: %v, actual testcase %v", expected, ts.Tests)
	}
}
