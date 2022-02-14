package shell

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

var (
	isDebug = false
)

func SetDebug(debug bool) {
	isDebug = debug
}

// Shell标准输出缓冲区
// 用于返回输出的内容
type shellStdBuffer struct {
	writer io.Writer
	buf    *bytes.Buffer
}

func newShellStdBuffer(writer io.Writer) *shellStdBuffer {
	return &shellStdBuffer{
		writer: writer,
		buf:    bytes.NewBuffer([]byte{}),
	}
}
func (s *shellStdBuffer) Write(p []byte) (n int, err error) {
	n, err = s.buf.Write(p)
	if s.writer != nil {
		n, err = s.writer.Write(p)
	}
	return n, err
}
func (s *shellStdBuffer) String() string {
	return string(s.buf.Bytes())
}

// 执行Shell命令
// 如果没有返回error,则命令执行成功，反之失败
// code返回命令执行返回的状态码,返回0表示执行成功
// output返回命令输出内容
func execCommand(command string, stdIn io.Reader, stdOut io.Writer,
	stdErr io.Writer, debug bool) (code int, output string, err error) {

	var status syscall.WaitStatus //执行状态
	//var output string             //输出内容
	var stdout *shellStdBuffer //标准输出
	var stderr *shellStdBuffer //标准错误输出

	if strings.TrimSpace(command) == "" {
		return 0, "", errors.New("no such command")
	}

	if debug {
		fmt.Println(fmt.Sprintf("[COMMAND]:\n%s\n%s",
			command, strings.Repeat("-", len(command))))
	}

	//var arr = string.Split(command, " ")
	var cmd *exec.Cmd
	if runtime.GOOS == "windows"{
		var arr = strings.Split(command, " ")
		 cmd = exec.Command(arr[0], arr[1:]...)
	}else {
		cmd = exec.Command("sh", "-c", command)
	}
	stdout = newShellStdBuffer(stdOut)
	stderr = newShellStdBuffer(stdErr)

	cmd.Stdout = stdout
	cmd.Stdin = stdIn
	cmd.Stderr = stderr

	err = cmd.Start()
	if err != nil {
		return 1, "", err
	}

	err = cmd.Wait()

	status = cmd.ProcessState.Sys().(syscall.WaitStatus)
	isSuccess := cmd.ProcessState.Success()
	if debug {
		fmt.Println(strings.Repeat("-", len(command)))
		if isSuccess {
			fmt.Println("[OK] Status:", status.ExitStatus(),
				" Used Time:", cmd.ProcessState.UserTime())
		} else {
			fmt.Println("[Fail] Status:", status.ExitStatus(),
				" Used Time:", cmd.ProcessState.UserTime())
		}
	}

	if isSuccess {
		output = stdout.String()
	} else {
		output = stderr.String()
	}

	return status.ExitStatus(), output, nil
}

// Run 执行Shell命令
// 如果没有返回error,则命令执行成功，反之失败
// code返回命令执行返回的状态码,返回0表示执行成功
// stdOutput:true 输出到os.StdOut, 错误输出到os.StdErr, 需要将正常结果输出到stderr中才能显示命令输出

func Run(command string,stdOutput bool) (code int, output string, err error) {
	//return execCommand(command, os.Stdin, os.Stdout, os.Stdin, isDebug)
	if stdOutput{
		return execCommand(command, os.Stderr, os.Stdout, os.Stderr, isDebug)
	}
	return execCommand(command, nil, nil, nil, isDebug)
}

// Brun (后台/静默)执行
// 仅仅执行命令，不需要捕获结果
func Brun(command string) (err error) {
	if strings.TrimSpace(command) == "" {
		return errors.New("no such command")
	}

	var arr = strings.Split(command, " ")
	var cmd = exec.Command(arr[0], arr[1:]...)
	err = cmd.Start()

	if isDebug {
		fmt.Print(fmt.Sprintf("[COMMAND]:\n%s\n%s",
			command, strings.Repeat("-", len(command))))
		if err != nil {
			fmt.Print("[ ERROR]:", err.Error())
		}
		fmt.Print("\n")
	}
	return err
}
