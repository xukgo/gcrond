package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/xukgo/gcrond/core"
)

func main() {
	projconf := new(core.ProjectConfig)
	task := new(core.TimerTaskConfig)
	task.DoShellPath = []string{"/opt/cmd/aa.sh"}
	task.IntervalStr = "1d"
	projconf.TimerTasks = []*core.TimerTaskConfig{task}
	ruleExecs := new(core.RuleExecConfig)
	ruleExecs.IntervalStr = "10s"
	ruleExecs.DoShellPath = []string{"/opt/cmd/bb.sh"}
	existConf := new(core.CheckExistConfig)
	existConf.ExecPath = "/bin/doa"
	existConf.IncludeCmd = []string{"111", "222"}
	existConf.ExcludeCmd = []string{"333", "444"}
	ruleExecs.CheckConfig = existConf
	projconf.RuleExecConfig = []*core.RuleExecConfig{ruleExecs}
	ymlstr, _ := yaml.Marshal(projconf)
	fmt.Println(ymlstr)

	core.ShellDo("/opt/newcc/start.sh")
	arr := core.GetProcessCmdLines("/opt/newcc/start.sh", []string{"newcc"}, nil)
	if len(arr) == 0 {

	}
}
