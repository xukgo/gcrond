package core

import (
	"bytes"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"github.com/xukgo/gcrond/compon"
	"github.com/xukgo/gcrond/logUtil"
	"go.uber.org/zap"
	"time"
)

type RuleExecConfig struct {
	Enable        bool              `yaml:"enable"`
	StartupDelay  int64             `yaml:"startupDelay"`
	Description   string            `yaml:"description"`
	IntervalStr   string            `yaml:"interval"`
	IntervalSec   int64             `yaml:"-"`
	Commands      []string          `yaml:"command"`
	CheckConfig   *CheckExistConfig `yaml:"check"`
	LastCheckUnix int64             `yaml:"-"`
}

type CheckExistConfig struct {
	ExecPath   string   `yaml:"execPath"`
	IncludeCmd []string `yaml:"includeCmd"`
	ExcludeCmd []string `yaml:"excludeCmd"`
}

func (this *RuleExecConfig) ToDescription() string {
	var bf bytes.Buffer
	bf.WriteString("解析启用规则任务\n")
	bf.WriteString(fmt.Sprintf("    Description=>%s\n", this.Description))
	bf.WriteString(fmt.Sprintf("    StartupDelay=>%d秒\n", this.StartupDelay))
	bf.WriteString(fmt.Sprintf("    Interval=>%s\n", this.IntervalStr))
	bf.WriteString(fmt.Sprintf("    Commands:\n"))
	for _, cmd := range this.Commands {
		bf.WriteString(fmt.Sprintf("        =>%s\n", cmd))
	}
	return bf.String()
}

func (this *RuleExecConfig) CheckParam() bool {
	if this.StartupDelay < 0 {
		logUtil.LoggerCommon.Error("RuleExecConfig StartupDelay is not valid")
		return false
	}
	if this.IntervalSec <= 0 {
		logUtil.LoggerCommon.Error("RuleExecConfig interval is not valid")
		return false
	}
	if len(this.Commands) == 0 {
		logUtil.LoggerCommon.Error("RuleExecConfig commands is not valid")
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
	if this.StartupDelay > 0 {
		duration, err := compon.GetSystemUptime()
		if err != nil {
			logUtil.LoggerCommon.Error("GetSystemUptime error", zap.Error(err))
			return
		}
		if duration.Seconds() < float64(this.StartupDelay) {
			return
		}
	}

	if time.Now().Unix()-this.LastCheckUnix < this.IntervalSec {
		return
	}

	cmds := GetProcessCmdLines(procInfos, this.CheckConfig.ExecPath, this.CheckConfig.IncludeCmd, this.CheckConfig.ExcludeCmd)
	if len(cmds) > 0 {
		return
	}

	procInfos, err := process.Processes()
	if err != nil {
		logUtil.LoggerCommon.Error("get process error", zap.Error(err))
		return
	}

	cmds = GetProcessCmdLines(procInfos, this.CheckConfig.ExecPath, this.CheckConfig.IncludeCmd, this.CheckConfig.ExcludeCmd)
	if len(cmds) > 0 {
		return
	}

	for _, cmd := range this.Commands {
		logUtil.LoggerCommon.Info("规则任务开始执行", zap.String("description", this.Description), zap.String("cmd", cmd))
		outStr, errStr, err := ExecCmdline(cmd)
		if err != nil {
			logUtil.LoggerCommon.Error("规则任务执行失败", zap.Error(err), zap.String("description", this.Description),
				zap.String("cmd", cmd), zap.String("stdout", outStr), zap.String("stderr", errStr))
			this.LastCheckUnix = time.Now().Unix()
			break
		} else {
			logUtil.LoggerCommon.Info("规则任务执行成功", zap.Error(err), zap.String("description", this.Description),
				zap.String("cmd", cmd), zap.String("stdout", outStr))
		}
	}

	this.LastCheckUnix = time.Now().Unix()
}
