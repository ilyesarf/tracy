package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type SysCall struct {
	Time string `json:"time"`
	OP   string `json:"op"`
}

type Trace struct {
	Binary   string    `json:"binary"`
	SysCalls []SysCall `json:"syscalls"`
}

func trace_bin(path string) {

	// Create the strace command
	cmd := exec.Command("blink", "-s", path)

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

}

func readLog() []string {
	logPath := "blink.log"

	file, err := os.Open(logPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var unparsed_syscalls []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		unparsed_syscalls = append(unparsed_syscalls, line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	os.Remove(logPath)
	return unparsed_syscalls
}

func parseLog(unparsed_syscalls []string) []SysCall {
	var syscalls []SysCall
	for _, line := range unparsed_syscalls {
		var syscall SysCall
		if strings.Contains(line, "strace.c") || strings.Contains(line, "exit_group") {
			splitted := strings.Split(line, " ")
			syscall.Time = strings.Split(splitted[0], ":blink")[0]
			syscall.OP = strings.Join(splitted[2:], " ")
			syscalls = append(syscalls, syscall)
		}
	}

	return syscalls
}

func main() {
	args := os.Args

	var path string
	if len(args) == 2 {
		path = args[1]
	} else {
		path = "tmp/a.out"
	}

	trace_bin(path)
	unparsed_syscalls := readLog()

	var trace Trace
	trace.Binary = path
	trace.SysCalls = parseLog(unparsed_syscalls)

	fmt.Println(trace.Binary)
	for _, syscall := range trace.SysCalls {
		fmt.Println(syscall)
	}
}
