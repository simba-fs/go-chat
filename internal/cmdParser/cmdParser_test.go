package cmdParser

import "testing"
import "errors"

func TestCmdType(t *testing.T) {
	c := []Cmd{
		{
			"/name",
			0,
			"/name",
			"show current name",
			func(raw string, cmds []string, exec FuncExec) (string, error) {
				return "simba", nil
			},
		},
		{
			"/room",
			0,
			"/room [room id]",
			"show or change current room",
			func(raw string, cmds []string, exec FuncExec) (string, error) {
				if len(cmds) > 2 {
					return cmds[1], nil
				}
				return "1234", nil
			},
		},
		{
			"/server",
			0,
			"/server [ip]",
			"show or change current server ip",
			func(raw string, cmds []string, exec FuncExec) (string, error) {
				if len(cmds) > 2 {
					return cmds[1], nil
				}
				return "127.0.0.1", nil
			},
		},
	}

	t.Log(c)
}

func TestCmdListTye(t *testing.T) {
	c := []Cmd{
		{
			"/name",
			0,
			"/name",
			"show current name",
			func(raw string, cmds []string, exec FuncExec) (string, error) {
				return "simba", nil
			},
		},
		{
			"/room",
			0,
			"/room [room id]",
			"show or change current room",
			func(raw string, cmds []string, exec FuncExec) (string, error) {
				if len(cmds) > 2 {
					return cmds[1], nil
				}
				return "1234", nil
			},
		},
		{
			"/server",
			0,
			"/server [ip]",
			"show or change current server ip",
			func(raw string, cmds []string, exec FuncExec) (string, error) {
				if len(cmds) > 2 {
					return cmds[1], nil
				}
				return "127.0.0.1", nil
			},
		},
	}

	cl := CmdList{
		c,
		"/help",
		Helper,
	}

	cl.Exec("/name")
	cl.Exec("/room 12")
	cl.Exec("/room")
	t.Log(cl.Exec("/help"))
}

func TestNew(t *testing.T){
	c := New("/connect [ip]", "connect to server", func(raw string, cmds []string, exec FuncExec)(string, error){
		if len(cmds) < 2 {
			return "", errors.New("length of arg not match")
		}
		return cmds[1], nil
	})

	t.Log(c)
}
