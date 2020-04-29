package core

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/xukgo/gcrond/logUtil"
	"go.uber.org/zap"
)

type ProjectConfig struct {
	TimerTasks     []*TimerTaskConfig `yaml:"TimerTask"`
	RuleExecConfig []*RuleExecConfig  `yaml:"RuleExec"`
}

func (this *ProjectConfig) FillWithYaml(data []byte) error {
	err := yaml.Unmarshal(data, this)
	if err != nil {
		logUtil.LoggerCommon.Error("ProjectConfig unmarshal yaml error", zap.Error(err))
		return err
	}
	for idx := range this.TimerTasks {
		this.TimerTasks[idx].IntervalSec, err = ParseInterval(this.TimerTasks[idx].IntervalStr)
		if err != nil {
			return err
		}
		if !this.TimerTasks[idx].CheckParam() {
			return fmt.Errorf("TimerTask校验参数未通过")
		}
	}
	for idx := range this.RuleExecConfig {
		this.RuleExecConfig[idx].IntervalSec, err = ParseInterval(this.RuleExecConfig[idx].IntervalStr)
		if err != nil {
			return err
		}
		if !this.RuleExecConfig[idx].CheckParam() {
			return fmt.Errorf("RuleExecConfig校验参数未通过")
		}
	}
	return nil
}
