package testcase

import (
	. "github.com/liquidz/shelltest/mock"
	"regexp"
	"strings"
)

const (
	EqualMethod  = "equals"
	RegexpMethod = "matches"
)

const (
	DefaultMethod = EqualMethod
)

var NoAutoAssertion = false
var commandRegexp = regexp.MustCompile(`^[^$]*\$\s+(.+)\s*$`)
var sectionRegexp = regexp.MustCompile(`^\[\s*(.+)\s*\]$`)
var newLineRegexp = regexp.MustCompile(`[\r\n]+`)
var multiLineRegexp = regexp.MustCompile(`\s+\\\s*[\r\n]+`)
var commentRegexp = regexp.MustCompile(`^\s*#\s*`)
var regexpRegexp = regexp.MustCompile(`^=~\s+(.+)\s*$`)
var requireRegexp = regexp.MustCompile(`^@require\s+(.+)\s*$`)

func getAutoAssertion(tc TestCase) TestCase {
	if tc.IsEmpty() || NoAutoAssertion || len(tc.Expected) > 0 {
		return TestCase{}
	}

	comment := tc.Comment
	if comment == "" {
		comment = tc.Command
	}

	return TestCase{
		Command:  "echo $?",
		Expected: Assertion{Method: EqualMethod, Text: "0"}.ToArray(),
		Comment:  comment,
	}
}

func Parse(filepath string) (TestSuite, error) {
	var (
		tc    TestCase
		ts    TestSuite
		match []string
	)

	b, err := ReadFile(filepath)
	if err != nil {
		return ts, err
	}

	lastComment := ""
	s := multiLineRegexp.ReplaceAllString(string(b), " ")
	for _, l := range newLineRegexp.Split(strings.TrimSpace(s), -1) {
		if strings.TrimSpace(l) == "" {
			continue
		}

		// comment
		if commentRegexp.MatchString(l) {
			lastComment = commentRegexp.ReplaceAllString(l, "")
			continue
		}

		// require
		if match = requireRegexp.FindStringSubmatch(l); len(match) == 2 {
			ts.Append(tc)
			ts.Append(getAutoAssertion(tc))
			tc = TestCase{}

			ts2, err := Parse(match[1])
			if err != nil {
				return ts, err
			}
			ts = ts.Merge(ts2)
			continue
		}

		// Command Line
		if match = commandRegexp.FindStringSubmatch(l); len(match) == 2 {
			ts.Append(tc)
			ts.Append(getAutoAssertion(tc))
			if match[1] == "exit" {
				tc = TestCase{}
			} else {
				tc = TestCase{Command: match[1], Comment: lastComment}
			}
			continue
		}

		// Regexp Assertion
		if match = regexpRegexp.FindStringSubmatch(l); len(match) == 2 {
			if _, err := regexp.Compile(match[1]); err != nil {
				return ts, err
			}
			tc.AppendAssertion(RegexpMethod, match[1])
			continue
		}

		tc.AppendAssertion(DefaultMethod, l)
	}

	ts.Append(tc)
	ts.Append(getAutoAssertion(tc))

	return ts, nil
}
