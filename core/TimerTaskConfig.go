package core

import "github.com/xukgo/gcrond/logUtil"

type TimerTaskConfig struct {
	IntervalStr string   `yaml:"interval"`
	DoShellPath []string `yaml:"shellPath"`
	IntervalSec int64    `yaml:"-"`
}

func (this *TimerTaskConfig) CheckParam() bool {
	if this.IntervalSec <= 0 {
		logUtil.LoggerCommon.Error("TimerTaskConfig interval is not valid")
		return false
	}
	if len(this.DoShellPath) == 0 {
		logUtil.LoggerCommon.Error("TimerTaskConfig shell path is not valid")
		return false
	}
	return true
}
