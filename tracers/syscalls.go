package tracers

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gorilla/websocket"
)

type SysCall struct {
	Time string `json:"time"`
	OP   string `json:"op"`
}

type Trace struct {
	Binary   string `json:"binary"`
	Args     []string
	SysCalls []SysCall `json:"syscalls"`
}

func (t *Trace) SendTrace(conn *websocket.Conn) {
	data, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	//fmt.Println(conn)
	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println(err)
		conn.Close()
	}
}

func (t *Trace) ParseLog(line string) {
	var syscall SysCall

	splitted := strings.Split(line, " ")
	syscall.Time = splitted[0]
	syscall.OP = strings.Join(splitted[1:], " ")
	t.SysCalls = append(t.SysCalls, syscall)
}

func (t *Trace) TraceBin() {
	cmdArgs := []string{"strace", "-tt", t.Binary}
	cmdArgs = append(cmdArgs, t.Args...)

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

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
		}

	}()

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}

}
