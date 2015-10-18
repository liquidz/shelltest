package testcase

import (
	"fmt"
	"regexp"
	"strings"
)

type TestCases []TestCase

type TestCase struct {
	Command  string
	Expected Assertions
	Comment  string
}

type Assertions []Assertion

type Assertion struct {
	Method string
	Text   string
}

func (a Assertion) ToArray() Assertions {
	return Assertions{a}
}

func (a *Assertion) Assert(s string) bool {
	switch a.Method {
	case EqualMethod:
		return (a.Text == s)
	case RegexpMethod:
		return regexp.MustCompile(a.Text).MatchString(s)
	}
	return false
}

func (as Assertions) IsExpected(s string) bool {
	r := regexp.MustCompile(`[\r\n]+`)
	outputs := r.Split(strings.TrimSpace(s), -1)
	outputLen := len(outputs)

	for i, expected := range as {
		if i >= outputLen {
			return false
		}

		actual := strings.TrimSpace(outputs[i])
		if !expected.Assert(actual) {
			return false
		}
	}
	return true
}

func (tc *TestCase) IsEmpty() bool {
	return tc.Command == "" && len(tc.Expected) == 0
}

func (tc *TestCase) AppendAssertion(method, text string) {
	if tc.IsEmpty() {
		return
	}
	tc.Expected = append(tc.Expected, Assertion{
		Method: strings.TrimSpace(method),
		Text:   strings.TrimSpace(text),
	})
}

func (tc *TestCase) String() string {
	return fmt.Sprintf("Command: %v, Assertion: %v", tc.Command, tc.Expected)
}

func (tcs TestCases) Length() int {
	n := 0
	for _, tc := range tcs {
		if len(tc.Expected) == 0 {
			continue
		}
		n++
	}
	return n
}
