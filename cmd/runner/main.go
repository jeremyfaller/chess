package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type prog struct {
	path    string
	prog    *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	scanner *bufio.Scanner
}

// NewProg makes a new prog interface.
//
// Note that the program has already been started.
func NewProg(name string) (*prog, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return nil, fmt.Errorf("error finding executable: %w", err)
	}
	cmd := exec.Command(path)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("error attaching to stdin: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error attaching to stout: %w", err)
	}
	scanner := bufio.NewScanner(stdout)
	go cmd.Run()
	return &prog{path: path, prog: cmd, stdin: stdin, stdout: stdout, scanner: scanner}, nil
}

// writeln writes to a prog.
//
// It doesn't wait for a response. If cmd doesn't contain a newline, one is added.
func (p *prog) writeln(cmd string) error {
	toWrite := cmd
	if !strings.HasSuffix(cmd, "\n") {
		toWrite += "\n"
	}
	n, err := p.stdin.Write([]byte(toWrite))
	if err != nil {
		return fmt.Errorf("error writing %q command: %w", strings.Trim(cmd, " \t\n\r"), err)
	}
	if n != len(toWrite) {
		return fmt.Errorf("didn't complete the write: %d != %d", n, len(toWrite))
	}
	return nil
}

// RunCommand runs a command to a program, waiting for an output.
func (p *prog) RunCommand(cmd string, output string) error {
	if err := p.writeln(cmd); err != nil {
		return err
	}

	for p.scanner.Scan() {
		str := strings.Trim(p.scanner.Text(), " \t\r\n")
		if strings.Contains(str, output) {
			return nil
		}
	}
	return nil
}

func run(p1, p2 *prog) error {
	for _, p := range []*prog{p1, p2} {
		if err := p.RunCommand("uci", "uciok"); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		log.Fatalf("usage: %v [flags] p1 p2", os.Args[0])
	}
	p1, err := NewProg(flag.Arg(0))
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	p2, err := NewProg(flag.Arg(1))
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := run(p1, p2); err != nil {
		log.Fatalf("error running: %v", err)
	}
}
