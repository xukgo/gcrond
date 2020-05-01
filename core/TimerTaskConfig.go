package core

import (
	"bytes"
	"fmt"
	"github.com/xukgo/gcrond/compon"
	"github.com/xukgo/gcrond/logUtil"
	"go.uber.org/zap"
	"time"
)

type TimerTaskConfig struct {
	Enable       bool     `yaml:"enable"`
	StartupDelay int64    `yaml:"startupDelay"`
	Description  string   `yaml:"description"`
	IntervalStr  string   `yaml:"interval"`
	Commands     []string `yaml:"command"`
	IntervalSec  int64    `yaml:"-"`
}

func (this *TimerTaskConfig) ToDescription() string {
	var bf bytes.Buffer
	bf.WriteString("解析启用定时任务\n")
	bf.WriteString(fmt.Sprintf("    Description=>%s\n", this.Description))
	bf.WriteString(fmt.Sprintf("    StartupDelay=>%d秒\n", this.StartupDelay))
	bf.WriteString(fmt.Sprintf("    Interval=>%s\n", this.IntervalStr))
	bf.WriteString(fmt.Sprintf("    Commands:\n"))
	for _, cmd := range this.Commands {
		bf.WriteString(fmt.Sprintf("        =>%s\n", cmd))
	}
	return bf.String()
}
func (this *TimerTaskConfig) CheckParam() bool {
	if this.StartupDelay < 0 {
		logUtil.LoggerCommon.Error("TimerTaskConfig StartupDelay is not valid")
		return false
	}
	if this.IntervalSec <= 0 {
		logUtil.LoggerCommon.Error("TimerTaskConfig interval is not valid")
		return false
	}
	if len(this.Commands) == 0 {
		logUtil.LoggerCommon.Error("TimerTaskConfig commands is not valid")
		return false
	}
	return true
}

func (this *TimerTaskConfig) Start() {
	if !this.Enable {
		return
	}

	go func() {
		logUtil.LoggerCommon.Info("定时任务启动", zap.String("description", this.Description))
		for {
			time.Sleep(time.Second)

			if this.StartupDelay > 0 {
				duration, err := compon.GetSystemUptime()
				if err != nil {
					logUtil.LoggerCommon.Error("GetSystemUptime error", zap.Error(err))
					return
				}
				if duration.Seconds() >= float64(this.StartupDelay) {
					break
				}
			} else {
				break
			}
		}
		for {
			for _, cmd := range this.Commands {
				logUtil.LoggerCommon.Info("定时任务开始执行", zap.String("description", this.Description), zap.String("cmd", cmd))
				outStr, errStr, err := ExecCmdline(cmd)
				if err != nil {
					logUtil.LoggerCommon.Error("定时任务执行失败", zap.Error(err), zap.String("description", this.Description),
						zap.String("cmd", cmd), zap.String("stdout", outStr), zap.String("stderr", errStr))
					break
				} else {
					logUtil.LoggerCommon.Info("定时任务执行成功", zap.Error(err), zap.String("description", this.Description),
						zap.String("cmd", cmd), zap.String("stdout", outStr))
				}
			}
			time.Sleep(time.Second * time.Duration(this.IntervalSec))
		}
	}()
}
