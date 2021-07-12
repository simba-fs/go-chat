package cmdParser

import (
	"errors"
	"fmt"
	"strings"
)

// FuncExec can execute aother command in command definition object (Cmd)
type FuncExec func(exec string, arg ...interface{}) (string, error)

// Cmd stores one command
type Cmd struct {
	Cmd         string
	ArgLen      int64 // minimum required length of arguments, cmd not included
	Usage       string
	Description string
	Exec        func(raw string, cmds []string, exec FuncExec, arg ...interface{}) (string, error)
}

// Cmds stores a lot of commands
type CmdList struct {
	Cmds   []Cmd
	Help   string
	Helper func(cmds *CmdList) string
}

var ErrCommandNotFound = errors.New("command not found")
var ErrMissArg = errors.New("missing args")

// Exec execute commands in cmdList
func (cmdList *CmdList) Exec(raw string, arg ...interface{}) (string, error) {
	cmds := splitCmd(raw)
	if cmds[0] == cmdList.Help || cmds[0] == "" {
		return cmdList.Helper(cmdList), nil
	}

	var exec Cmd
	flag := false

	for _, i := range cmdList.Cmds {
		if i.Cmd != cmds[0] {
			continue
		}

		if int64(len(cmds)-1) < i.ArgLen {
			return "", ErrMissArg
		}

		exec = i
		flag = true
	}

	if !flag {
		return "", ErrCommandNotFound
	}

	return exec.Exec(raw, cmds, cmdList.Exec, arg...)

}

// Helper is default helper function, you can custom it
func Helper(cmds *CmdList) string {
	helpMsg := ""
	for _, i := range cmds.Cmds {
		helpMsg += fmt.Sprintf("%s   %s\n", i.Usage, i.Description)
	}
	return helpMsg
}

// New return a Cmd object
// example:
//     cmdParser.New("/connect <ip> [port=8888]", "connect to server", func(raw string, cmds []string, exec cmdParser.FuncExec)(string, error){
//          server := cmd[1]
//          ip := 8888
//          if len(cmd) >= 3 {
//              ip = cmd[2]
//          }
//          fmt.Printf("connected to server %s:%s", server, ip)
//
//          retrurn "", nil
//     })
func New(usage string, description string, exec func(raw string, cmds []string, exec FuncExec, arg ...interface{}) (string, error)) Cmd {
	cmd := splitCmd(usage)

	c := Cmd{
		Cmd:         cmd[0],
		Usage:       usage,
		Description: description,
		Exec:        exec,
	}

	argLen := int64(0)
	for _, i := range cmd[1:] {
		head := i[0]
		tail := i[len(i)-1]
		if head == '<' && tail == '>' {
			argLen++
		}
	}

	c.ArgLen = argLen

	return c
}

// splitCmd split a string into string slice, spliting by space, removing space
func splitCmd(raw string) []string {
	if len(raw) <= 0 {
		return []string{""}
	}

	cmd := []string{}
	for _, val := range strings.Split(raw, " ") {
		if val != "" {
			cmd = append(cmd, val)
		}
	}

	if len(cmd) == 0 {
		cmd = append(cmd, "")
	}
	return cmd
}
