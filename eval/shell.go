package eval

import (
	"bufio"
	"bytes"
	"fmt"
	. "github.com/liquidz/shelltest/color"
	. "github.com/liquidz/shelltest/debug"
	"io"
	"os/exec"
	"strings"
)

const (
	Delimiter = "shelltest_asd8f9s8df0as8df90a8sfd098a0s8"
)

func inputLoop(stdin io.WriteCloser, inCh chan string) {
	for {
		s := <-inCh
		if s == "exit" {
			io.WriteString(stdin, "exit\n")
			break
		} else {
			DebugPrint("shell", "run [%v]", s)
			io.WriteString(stdin, fmt.Sprintf("echo %v\n", Delimiter))
			io.WriteString(stdin, fmt.Sprintf("(exit $RT); %v; RT=$?\n", s))
		}
	}
}

func scanLoop(scanner *bufio.Scanner, outCh chan string, outTermCh chan bool) {
	first := true
	buf := new(bytes.Buffer)

	for scanner.Scan() {
		text := scanner.Text()
		if DeleteColor(text) == Delimiter {
			if first {
				first = false
			} else {
				text := strings.TrimSpace(buf.String())
				DebugPrint("shell", "out <- [%v]", text)
				outCh <- text
				buf = new(bytes.Buffer)
			}
		} else {
			DebugPrint("shell", "read [%v]", text)
			buf.WriteString(text + "\n")
		}
	}
	text := strings.TrimSpace(buf.String())
	DebugPrint("shell", "out <- [%v]", text)
	if text != "" {
		outCh <- text
	}
	outTermCh <- true
}

func startShell(shell string, initCommands ...string) (chan string, chan string, chan bool, error) {
	inCh := make(chan string)
	outCh := make(chan string)
	outTermCh := make(chan bool)
	termCh := make(chan bool)

	cmd := exec.Command("bash")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, nil, nil, err
	}

	io.WriteString(stdin, "\n")
	for _, c := range initCommands {
		io.WriteString(stdin, c+"\n")
	}
	outScanner := bufio.NewScanner(stdout)

	go inputLoop(stdin, inCh)
	go scanLoop(outScanner, outCh, outTermCh)
	go func() {
		cmd.Wait()
		<-outTermCh
		termCh <- true
	}()

	return inCh, outCh, termCh, nil
}
