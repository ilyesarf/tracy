package tracers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
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

func (t *Trace) sendTrace() {
	endp := "http://localhost:1337/sendTrace"
	body, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	var r *http.Request
	r, err = http.NewRequest("POST", endp, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
}

func (t *Trace) ParseLog(line string) {
	var syscall SysCall
	if strings.Contains(line, "strace.c") || strings.Contains(line, "exit_group") {
		splitted := strings.Split(line, " ")
		syscall.Time = strings.Split(splitted[0], ":blink")[0]
		syscall.OP = strings.Join(splitted[2:], " ")
		t.SysCalls = append(t.SysCalls, syscall)
	}

}

func (t *Trace) TraceBin() {
	cmd := exec.Command("blink", "-s", "-e", t.Binary)

	// Create pipes for standard input/output/error
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			output := scanner.Text()
			t.ParseLog(output)
			t.sendTrace()
		}

	}()

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}

}
