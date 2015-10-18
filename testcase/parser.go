package testcase

import (
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

var NoAutoAssertion = false
var commandRegexp = regexp.MustCompile(`^[^$]*\$\s+(.+)\s*$`)
var sectionRegexp = regexp.MustCompile(`^\[\s*(.+)\s*\]$`)
var newLineRegexp = regexp.MustCompile(`[\r\n]+`)
var multiLineRegexp = regexp.MustCompile(`\s+\\\s*[\r\n]+`)
var commentRegexp = regexp.MustCompile(`^\s*#\s*`)
var regexpRegexp = regexp.MustCompile(`^=~\s+(.+)\s*$`)

func addAutoAssertion(tc TestCase) TestCase {
	if tc.IsEmpty() || NoAutoAssertion || len(tc.Expected) > 0 {
		return tc
	}

	return TestCase{
		Command:  tc.Command + " > /dev/null 2>&1; echo $?",
		Expected: Assertion{Method: EqualMethod, Text: "0"}.ToArray(),
		Comment:  tc.Comment,
	}
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
			ts.Append(section, addAutoAssertion(tc))
			tc = TestCase{}
			section = match[1]
			continue
		}

		match = commandRegexp.FindStringSubmatch(l)
		if len(match) == 2 {
			// Command Line
			ts.Append(section, addAutoAssertion(tc))
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

	ts.Append(section, addAutoAssertion(tc))

	return ts, nil
}
