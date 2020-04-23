package main

import (
	"github.com/xukgo/gcrond/core"
)

func main() {
	core.ShellDo("/opt/newcc/start.sh")
	arr := core.GetProcessCmdLines([]string{"newcc"}, nil)
	if len(arr) == 0 {

	}
}
