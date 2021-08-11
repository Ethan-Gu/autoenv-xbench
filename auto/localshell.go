package auto

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"time"
)

// Execute command with result back
func ExecCommand(s string) string {
	cmd := exec.Command("/bin/bash", "-c", s)

	var out bytes.Buffer
	cmd.Stdout = &out

	//Wait until cmd is executed
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Execute command err: %s", s)
	}
	return out.String()
}

// Execute command with interaction
func ExecCommandByLine(commandName string, params []string) bool {

	cmd := exec.Command(commandName, params...)
	fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}

	cmd.Start()
	//Create a stream to read by line
	reader := bufio.NewReader(stdout)

	//Read each line of the stream
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}

	cmd.Wait()
	return true
}

// Execute command without result
func ExecCommandNoResult(command string) {

	cmd := exec.Command("/bin/bash", "-c", command)

	err := cmd.Start()

	if err != nil {
		fmt.Printf("%v: exec command:%v error:%v\n", time.Now(), command, err)
	}
	fmt.Printf("Waiting for command:%v to finish...\n", command)

	//Wait for the result of fork
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("%v: Command finished with error: %v\n", time.Now(), err)
	}

}
