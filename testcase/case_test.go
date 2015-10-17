package testcase

import (
	"reflect"
	"testing"
)

func TestAssertionAssertWithEqualMethod(t *testing.T) {
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

func TestAssertionAssertWithRegexpMethod(t *testing.T) {
	var a Assertion

	a = Assertion{Method: RegexpMethod, Text: "foo"}
	tests := map[string]bool{
		"":       false,
		"foo":    true,
		"bar":    false,
		"foobar": true,
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
