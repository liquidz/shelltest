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
	ColorRed          = 31
	ColorGreen        = 32
	ColorYellow       = 33
)

type CLI struct {
	outStream, errStream io.Writer
	suite                TestSuite
	nocolor              bool
}

func (cli *CLI) colorize(color int, s string) string {
	if cli.nocolor {
		return s
	} else {
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, s)
	}
}

func (cli *CLI) out(color int, format string, a ...interface{}) {
	format = cli.colorize(color, format+"\n")
	fmt.Fprintf(cli.outStream, format, a...)
}

func (cli *CLI) err(format string, a ...interface{}) {
	format = cli.colorize(ColorRed, format+"\n")
	fmt.Fprintf(cli.errStream, format, a...)
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
	flags.BoolVar(&cli.nocolor, "nocolor", false, "no color")
	flags.BoolVar(&flagVersion, "v", false, "Print version information and quit.")

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
		cli.out(ColorGreen, "success to parse")
		fmt.Fprintf(cli.outStream, "%v\n", cli.suite.String())
		return ExitCodeOK
	}

	errs := Evaluate(shell, cli.suite, func(_ TestCase, err error) {
		if err == nil {
			if !ShellTestDebugMode {
				fmt.Fprint(cli.outStream, ".")
			}
		} else {
			cli.out(ColorRed, "\n%v\n", err)
			if !ShellTestDebugMode {
				fmt.Fprint(cli.outStream, "x")
			}
		}
	})

	color := ColorGreen
	if len(errs) > 0 {
		color = ColorRed
	}

	cli.out(color, "\n\n%v tests, %v failures", len(cli.suite.Tests), len(errs))

	if len(errs) > 0 {
		return ExitCodeError
	}

	return ExitCodeOK
}
