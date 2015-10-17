package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	EqualMethod  = "equals"
	RegexpMethod = "matches"
)

const (
	DefaultSection = "test"
	DefaultMethod  = EqualMethod
)

var commandRegexp = regexp.MustCompile(`^[^$]*\$\s+(.+)\s*$`)
var sectionRegexp = regexp.MustCompile(`^\[\s*(.+)\s*\]$`)
var newLineRegexp = regexp.MustCompile(`[\r\n]+`)
var multiLineRegexp = regexp.MustCompile(`\s+\\\s*[\r\n]+`)
var commentRegexp = regexp.MustCompile(`^\s*#\s*`)
var regexpRegexp = regexp.MustCompile(`^=~\s+(.+)\s*$`)

type TestSuite struct {
	Before TestCases
	After  TestCases
	Tests  TestCases
}

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

func (ts *TestSuite) Append(section string, test TestCase) error {
	if test.IsEmpty() {
		return nil
	}

	switch section {
	case "before":
		ts.Before = append(ts.Before, test)
		return nil
	case "after":
		ts.After = append(ts.After, test)
		return nil
	case "test":
		ts.Tests = append(ts.Tests, test)
		return nil
	}
	return errors.New(fmt.Sprintf("unknown section: %v", section))
}

func (tc *TestCase) String() string {
	return fmt.Sprintf("Command: %v, Assertion: %v", tc.Command, tc.Expected)
}

func (ts *TestSuite) String() string {
	s := fmt.Sprintf(`
[before]
  %v
[tests]
  %v
[after]
  %v
	`, ts.Before, ts.Tests, ts.After)

	return strings.TrimSpace(s)
}

func Parse(s string) (TestSuite, error) {
	var (
		tc    TestCase
		ts    TestSuite
		match []string
	)

	lastComment := ""

	section := DefaultSection
	s = multiLineRegexp.ReplaceAllString(s, " ")

	for _, l := range newLineRegexp.Split(strings.TrimSpace(s), -1) {
		if strings.TrimSpace(l) == "" {
			continue
		}

		if commentRegexp.MatchString(l) {
			lastComment = commentRegexp.ReplaceAllString(l, "")
			continue
		}

		match = sectionRegexp.FindStringSubmatch(l)
		if len(match) == 2 {
			ts.Append(section, tc)
			tc = TestCase{}
			section = match[1]
			continue
		}

		match = commandRegexp.FindStringSubmatch(l)
		if len(match) == 2 {
			// Command Line
			ts.Append(section, tc)
			if match[1] == "exit" {
				tc = TestCase{}
			} else {
				tc = TestCase{Command: match[1], Comment: lastComment}
			}
		} else {
			// Result Line
			match = regexpRegexp.FindStringSubmatch(l)
			if len(match) == 2 {
				if _, err := regexp.Compile(match[1]); err != nil {
					return ts, err
				}
				tc.AppendAssertion(RegexpMethod, match[1])
			} else {
				tc.AppendAssertion(DefaultMethod, l)
			}
		}
	}

	ts.Append(section, tc)

	return ts, nil
}
