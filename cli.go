package main

import (
	"flag"
	"fmt"
	"io"
)

const (
	DefaultShell      = "bash"
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

type CLI struct {
	outStream, errStream io.Writer
	suite                TestSuite
}

func (cli *CLI) out(format string, a ...interface{}) {
	fmt.Fprintf(cli.outStream, format+"\n", a...)
}

func (cli *CLI) err(format string, a ...interface{}) {
	fmt.Fprintf(cli.errStream, format+"\n", a...)
}

func (cli *CLI) Run(args []string) int {
	var (
		flagLint    bool
		flagVersion bool
		shell       string
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&shell, "s", DefaultShell, "shell")
	flags.BoolVar(&flagLint, "l", false, "lint")
	flags.BoolVar(&ShellTestDebugMode, "d", false, "Debug mode")
	flags.BoolVar(&flagVersion, "v", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	if flagVersion {
		cli.err("%s version %s", Name, Version)
		return ExitCodeOK
	}

	args = flags.Args()
	if len(args) == 0 {
		cli.err("no argument")
		return ExitCodeError
	}

	b, err := ReadFile(args[0])
	if err != nil {
		cli.err("file is not found: %v", args[0])
		return ExitCodeError
	}

	cli.suite, err = Parse(string(b))
	if err != nil {
		cli.err("failed to parse: %v", err)
	}

	if flagLint {
		cli.out("success to parse\n%v", cli.suite.String())
		return ExitCodeOK
	}

	errs := Evaluate(shell, cli.suite, func(_ TestCase, err error) {
		if err == nil {
			fmt.Fprint(cli.outStream, ".")
		} else {
			cli.err("\n%v\n", err)
			fmt.Fprint(cli.errStream, "x")
		}
	})

	cli.out("\n\n%v tests, %v failures", len(cli.suite.Tests), len(errs))

	if len(errs) > 0 {
		return ExitCodeError
	}

	return ExitCodeOK
}
