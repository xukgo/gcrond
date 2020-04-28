package core

import (
	"github.com/go-yaml/yaml"
	"github.com/xukgo/gcrond/logUtil"
	"go.uber.org/zap"
	"time"
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

func (this *RuleExecConfig) FillWithYaml(data []byte) error {
	err := yaml.Unmarshal(data, this)
	if err != nil {
		logUtil.LoggerCommon.Error("RuleExecConfig unmarshal yaml error", zap.Error(err))
		return err
	}

	this.IntervalSec, err = ParseInterval(this.IntervalStr)
	if err != nil {
		return err
	}
	return nil
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

func (this *RuleExecConfig) Start() {
	go func() {
		for {
			this.checkAndDo()
			time.Sleep(time.Second * time.Duration(this.IntervalSec))
		}
	}()
}

func (this *RuleExecConfig) checkAndDo() {
	cmds := GetProcessCmdLines(this.CheckConfig.ExecPath, this.CheckConfig.IncludeCmd, this.CheckConfig.ExcludeCmd)
	if len(cmds) == 0 {
		return
	}

	for _, path := range this.DoShellPath {
		ShellDo(path)
	}
}
