package main

import (
	"errors"
	"flag"
	"fmt"
	. "github.com/liquidz/shelltest/color"
	. "github.com/liquidz/shelltest/debug"
	. "github.com/liquidz/shelltest/eval"
	. "github.com/liquidz/shelltest/formatter"
	. "github.com/liquidz/shelltest/testcase"
	"io"
	"os"
	"strings"
)

const (
	DefaultShell      = "bash"
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

type Env struct {
	Map map[string]string
}

func (e *Env) String() string {
	return fmt.Sprintf("%v", e.Map)
}

func (e *Env) Set(value string) error {
	kv := strings.SplitN(value, "=", 2)
	if len(kv) != 2 {
		return errors.New(fmt.Sprintf("fixme"))
	}

	e.Map[kv[0]] = kv[1]
	return nil
}

type CLI struct {
	outStream, errStream io.Writer
	suite                TestSuite
}

func (cli *CLI) out(format string, a ...interface{}) {
	fmt.Fprintf(cli.outStream, format+"\n", a...)
}

func (cli *CLI) err(format string, a ...interface{}) {
	format = RedStr(format + "\n")
	fmt.Fprintf(cli.errStream, format, a...)
}

func (cli *CLI) Run(args []string) int {
	var (
		flagLint    bool
		flagVersion bool
		shell       string
		fmtr        string
	)

	cwd, _ := os.Getwd()
	env := Env{map[string]string{"CWD": cwd}}

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&shell, "s", DefaultShell, "shell")
	flags.StringVar(&fmtr, "f", "default", "formatter")
	flags.BoolVar(&flagLint, "l", false, "lint")
	flags.BoolVar(&ShellTestDebugMode, "d", false, "Debug mode")
	flags.BoolVar(&NoColor, "nocolor", false, "no color")
	flags.BoolVar(&NoAutoAssertion, "noautoassert", false, "no auto assert")
	flags.BoolVar(&flagVersion, "v", false, "Print version information and quit.")
	flags.Var(&env, "E", "environmental variable")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	if flagVersion {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	args = flags.Args()
	if len(args) == 0 {
		cli.err("no argument")
		return ExitCodeError
	}

	// parse files
	for _, filepath := range args {
		ts, err := Parse(filepath)
		if err != nil {
			cli.err("failed to parse: %v", err)
			return ExitCodeError
		}
		cli.suite = cli.suite.Merge(ts)
	}
	cli.suite.EnvMap = env.Map

	if flagLint {
		cli.out(GreenStr("success to parse"))
		fmt.Fprintf(cli.outStream, "%v\n", cli.suite.String())
		return ExitCodeOK
	}

	formatter := SelectFormatter(fmtr)

	if s := formatter.Setup(cli.suite); s != "" {
		fmt.Fprintf(cli.outStream, s)
	}

	errs := Evaluate(shell, cli.suite, func(no int, tc TestCase, err error) {
		if s := formatter.Result(no, tc, err); s != "" {
			fmt.Fprintf(cli.outStream, s)
		}
	})

	if s := formatter.TearDown(cli.suite, errs); s != "" {
		fmt.Fprintf(cli.outStream, s)
	}

	if len(errs) > 0 {
		return ExitCodeError
	}
	return ExitCodeOK
}
