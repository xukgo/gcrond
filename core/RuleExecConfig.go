package core

import (
	"github.com/shirou/gopsutil/process"
	"github.com/xukgo/gcrond/logUtil"
	"go.uber.org/zap"
)

type RuleExecConfig struct {
	IntervalStr string            `yaml:"interval"`
	IntervalSec int64             `yaml:"-"`
	DoShellPath []string          `yaml:"shellPath"`
	CheckConfig *CheckExistConfig `yaml:"check"`
}

type CheckExistConfig struct {
	ExecPath   string   `yaml:"execPath"`
	IncludeCmd []string `yaml:"includeCmd"`
	ExcludeCmd []string `yaml:"excludeCmd"`
}

func (this *RuleExecConfig) CheckParam() bool {
	if this.IntervalSec <= 0 {
		logUtil.LoggerCommon.Error("RuleExecConfig interval is not valid")
		return false
	}
	if len(this.DoShellPath) == 0 {
		logUtil.LoggerCommon.Error("RuleExecConfig shell path is not valid")
		return false
	}
	if this.CheckConfig == nil {
		logUtil.LoggerCommon.Error("RuleExecConfig check exist config is not valid")
		return false
	}
	if len(this.CheckConfig.ExecPath) == 0 && len(this.CheckConfig.IncludeCmd) == 0 {
		logUtil.LoggerCommon.Error("RuleExecConfig check exist config not allow include param and exec path both are empty")
		return false
	}
	return true
}

func (this *RuleExecConfig) CheckAndDo(procInfos []*process.Process) {
	cmds := GetProcessCmdLines(procInfos, this.CheckConfig.ExecPath, this.CheckConfig.IncludeCmd, this.CheckConfig.ExcludeCmd)
	if len(cmds) == 0 {
		return
	}

	procInfos, err := process.Processes()
	if err != nil {
		logUtil.LoggerCommon.Error("get process error", zap.Error(err))
		return
	}

	cmds = GetProcessCmdLines(procInfos, this.CheckConfig.ExecPath, this.CheckConfig.IncludeCmd, this.CheckConfig.ExcludeCmd)
	if len(cmds) == 0 {
		return
	}

	for _, path := range this.DoShellPath {
		ShellDo(path)
	}
}
