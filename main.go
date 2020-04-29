package main

import (
	"github.com/shirou/gopsutil/process"
	"github.com/xukgo/gcrond/core"
	"github.com/xukgo/gcrond/logUtil"
	"github.com/xukgo/gsaber/utils/fileUtil"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

func main() {
	filePath := fileUtil.GetAbsUrl("conf/crond.yml")
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		logUtil.LoggerCommon.Error("read file error", zap.Error(err))
	}
	conf := new(core.ProjectConfig)
	err = conf.FillWithYaml(content)
	if err != nil {
		logUtil.LoggerCommon.Error("ProjectConfig unmarshal error", zap.Error(err))
	}

	if len(conf.RuleExecConfig) > 0 {
		for {
			procInfos, err := process.Processes()
			if err != nil {
				logUtil.LoggerCommon.Error("get process error", zap.Error(err))
				time.Sleep(time.Second)
				continue
			}

			for _, ruleExec := range conf.RuleExecConfig {
				ruleExec.CheckAndDo(procInfos)
			}
			time.Sleep(time.Second)
		}
	}

	for {
		time.Sleep(time.Hour)
	}
}
