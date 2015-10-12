package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	EqualMethod = "equals"
)

const (
	DefaultSection  = "test"
	DefaultMethod   = EqualMethod
	CommandRegexp   = `^[^$]*\$\s*(.+)\s*$`
	SectionRegexp   = `^\[\s*(.+)\s*\]$`
	NewLineRegexp   = `[\r\n]+`
	MultiLineRegexp = `\s+\\\s*[\r\n]+`
	CommentRegexp   = `^\s*#`
)

type TestSuite struct {
	Before TestCases
	After  TestCases
	Tests  TestCases
}

type TestCases []TestCase

type TestCase struct {
	Command  string
	Expected Assertions
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
	result := false
	switch a.Method {
	case "equals":
		result = (a.Text == s)
	}
	return result
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

	section := DefaultSection
	mr := regexp.MustCompile(MultiLineRegexp)
	sr := regexp.MustCompile(SectionRegexp)
	cr := regexp.MustCompile(CommandRegexp)
	nr := regexp.MustCompile(NewLineRegexp)
	comment := regexp.MustCompile(CommentRegexp)

	s = mr.ReplaceAllString(s, " ")

	for _, l := range nr.Split(strings.TrimSpace(s), -1) {
		if strings.TrimSpace(l) == "" {
			continue
		}

		if comment.MatchString(l) {
			continue
		}

		match = sr.FindStringSubmatch(l)
		if len(match) == 2 {
			ts.Append(section, tc)
			tc = TestCase{}
			section = match[1]
			continue
		}

		match = cr.FindStringSubmatch(l)
		if len(match) == 2 {
			ts.Append(section, tc)
			if match[1] == "exit" {
				tc = TestCase{}
			} else {
				tc = TestCase{Command: match[1]}
			}
		} else {
			tc.AppendAssertion(DefaultMethod, l)
		}
	}

	ts.Append(section, tc)

	return ts, nil
}
