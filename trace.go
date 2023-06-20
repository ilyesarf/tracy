package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

type SysCall struct {
	Time string
	OP   []string
}

type Trace struct {
	Binary   string
	SysCalls []SysCall
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

	var syscalls []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		syscalls = append(syscalls, line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	os.Remove(logPath)
	return syscalls
}

func main() {
	args := os.Args

	var path string
	if len(args) == 2 {
		path = args[1]
		fmt.Println(path)
	} else {
		path = "tmp/a.out"
	}

	trace_bin(path)

	syscalls := readLog()
	for _, syscall := range syscalls {
		fmt.Println(syscall)
	}
}
