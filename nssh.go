package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
)

var (
	addPrefix   bool
	version     string
	handleStdin bool
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
	flag.BoolVar(&handleStdin, "i", false, "handle STDIN")
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
	stdinChs := make([]chan ([]byte), len(targets))
	if handleStdin {
		for i, _ := range targets {
			stdinChs[i] = make(chan []byte, 256)
		}
		go func() {
			wg.Add(1)
			processStdin(stdinChs)
			wg.Done()
		}()
	}
	for i, host := range targets {
		wg.Add(1)
		go func(h string, ch chan []byte) {
			remoteCommand(h, ch, command)
			wg.Done()
		}(host, stdinChs[i])
	}
	wg.Wait()
}

func processStdin(chs []chan []byte) {
	buf := make([]byte, 1024)
	for {
		n, err := io.ReadAtLeast(os.Stdin, buf, 1)
		if err != nil {
			if err != io.EOF {
				log.Println("[error]", err)
			}
			break
		}
		for _, ch := range chs {
			ch <- buf[0:n]
		}
	}
	for _, ch := range chs {
		close(ch)
	}
}

func remoteCommand(host string, src chan []byte, command []string) {
	args := []string{}
	if src == nil {
		// man ssh(1): Multiple -t options force tty allocation, even if ssh has no local tty.
		args = append(args, "-t", "-t")
	}
	args = append(args, host)
	args = append(args, command...)

	cmd := exec.Command("ssh", args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "[error]", host, err)
		return
	}
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
	if src != nil {
		go writeInput(src, stdin, host)
	}
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

func writeInput(src chan []byte, dest io.WriteCloser, host string) {
	for {
		b, more := <-src
		if more {
			_, err := dest.Write(b)
			if err != nil {
				log.Println("[error]", host, err)
				dest.Close()
				break
			}
		} else {
			dest.Close()
			break
		}
	}
}
