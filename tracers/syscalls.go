package tracers

import (
	"bufio"
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

func (t *Trace) ExecBin() {

	cmd := exec.Command("blink", "-s", t.Binary)

	// Create pipes for standard input/output/error
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

}

func ReadLog() []string {
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

func ParseLog(unparsed_syscalls []string) []SysCall {
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

func (t *Trace) TraceBin() {

	t.ExecBin()

	unparsed_syscalls := ReadLog()
	t.SysCalls = ParseLog(unparsed_syscalls)

}
