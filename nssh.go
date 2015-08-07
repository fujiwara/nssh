package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

var (
	addPrefix bool
	version   string
)

type strslice []string

func (s *strslice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *strslice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func main() {
	var targets strslice
	var showVersion bool
	flag.Var(&targets, "t", "target hostname")
	flag.BoolVar(&addPrefix, "p", false, "add hostname to line prefix")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.Parse()
	if showVersion {
		fmt.Println("version:", version)
		return
	}

	command := flag.Args()
	if len(command) < 1 {
		flag.PrintDefaults()
		return
	}
	wg := &sync.WaitGroup{}
	for _, host := range targets {
		wg.Add(1)
		go func(h string) {
			remoteCommand(h, command)
			wg.Done()
		}(host)
	}
	wg.Wait()
}

func remoteCommand(host string, command []string) {
	args := []string{"-t", "-t"} // man ssh(1): Multiple -t options force tty allocation, even if ssh has no local tty.
	args = append(args, host)
	args = append(args, command...)

	cmd := exec.Command("ssh", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "[error]", host, err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "[error]", host, err)
		return
	}
	cmd.Start()
	go scanLines(stderr, os.Stderr, host)
	go scanLines(stdout, os.Stdout, host)
	cmd.Wait()
}

func scanLines(src io.ReadCloser, dest io.Writer, prefix string) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		if addPrefix {
			fmt.Fprintln(dest, prefix, scanner.Text())
		} else {
			fmt.Fprintln(dest, scanner.Text())
		}
	}
}
