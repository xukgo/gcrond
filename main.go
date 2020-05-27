package main

import (
	"flag"
	"fmt"
	_ "github.com/spf13/pflag"
	"github.com/xukgo/gcrond/compon/procUnique"
	"github.com/xukgo/gcrond/core"
	"github.com/xukgo/gcrond/logUtil"
	"github.com/xukgo/gcrond/psutil"
	"github.com/xukgo/gsaber/utils/fileUtil"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	//_ "net/http/pprof"
	"os"
	"time"
)

var (
	fhelp bool

	fversion bool
	fdaemon  bool

	fsearchKey string
)

func initFlag() {
	flag.BoolVar(&fhelp, "h", false, "this help")

	flag.BoolVar(&fversion, "v", false, "show version and exit")
	flag.BoolVar(&fdaemon, "d", false, "run as service")

	// 注意 `signal`。默认是 -s string，有了 `signal` 之后，变为 -s signal
	flag.StringVar(&fsearchKey, "e", "", "search key word and list process info")

	// 改变默认的 Usage，flag包中的Usage 其实是一个函数类型。这里是覆盖默认函数实现，具体见后面Usage部分的分析
	flag.Usage = usage
}
func usage() {
	fmt.Fprintf(os.Stderr, `gcrond version: gcrond/1.0.0
Options:
`)
	flag.PrintDefaults()
}

func main() {
	initFlag()
	flag.Parse()

	if fhelp {
		flag.Usage()
		return
	}

	if fversion {
		fmt.Fprintf(os.Stdout, "gcrond 1.0.0\n")
		return
	}

	if !fdaemon {
		if len(fsearchKey) == 0 {
			fmt.Fprintf(os.Stdout, "search key cannot be empty\n")
			return
		}

		procs, err := psutil.FilterGetProcCmdInfos()
		if err != nil {
			fmt.Fprintf(os.Stdout, "get process error:%s\n", err.Error())
			return
		}
		procInfos := core.GetProcess(procs, "", []string{fsearchKey}, nil)
		if len(procInfos) > 0 {
			for _, info := range procInfos {
				exe := info.Exe
				cmdline := info.Cmdline
				fmt.Fprintf(os.Stdout, "%s  =>  %s\n", exe, cmdline)
			}
			return
		}
	}

	var err error
	var procLocker = procUnique.NewLocker("gcrond_hms")
	err = procLocker.Lock()
	if err != nil {
		log.Println("应用不允许重复运行")
		os.Exit(-1)
	}

	defer procLocker.Unlock()

	logUtil.InitLogger()

	filePath := fileUtil.GetAbsUrl("conf/excludePrefix.xml")
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		logUtil.LoggerCommon.Error("read file error", zap.Error(err))
		return
	}
	excludeConf := new(core.ExcludeCommandXmlRoot)
	err = excludeConf.FillWithXml(string(content))
	if err != nil {
		logUtil.LoggerCommon.Error("ExcludeCommand unmarshal error", zap.Error(err))
		return
	}

	filePath = fileUtil.GetAbsUrl("conf/crond.yml")
	content, err = ioutil.ReadFile(filePath)
	if err != nil {
		logUtil.LoggerCommon.Error("read file error", zap.Error(err))
		return
	}

	conf := new(core.ProjectConfig)
	err = conf.FillWithYaml(content)
	if err != nil {
		logUtil.LoggerCommon.Error("ProjectConfig unmarshal error", zap.Error(err))
		return
	}

	for _, ruleExec := range conf.RuleExecConfig {
		if !ruleExec.Enable {
			continue
		}
		logUtil.LoggerCommon.Info(ruleExec.ToDescription())
	}

	//go func() {
	//	logUtil.LoggerCommon.Info("StartProfWebService")
	//	err := http.ListenAndServe(":60044", nil)
	//	if err != nil {
	//		logUtil.LoggerCommon.Error("StartProfWebService error", zap.Error(err))
	//	}
	//}()

	go startTimers(conf.TimerTasks)

	if len(conf.RuleExecConfig) > 0 {
		enableRuleExecArr := make([]*core.RuleExecConfig, 0, 4)
		enableCount := 0
		for _, ruleExec := range conf.RuleExecConfig {
			if ruleExec.Enable {
				enableCount++
				enableRuleExecArr = append(enableRuleExecArr, ruleExec)
			}
		}
		if enableCount > 0 {
			holdAllEnable(enableRuleExecArr)
		}
	}

	for {
		time.Sleep(time.Hour)
	}
}

func startTimers(tasks []*core.TimerTaskConfig) {
	for _, task := range tasks {
		if !task.Enable {
			continue
		}
		logUtil.LoggerCommon.Info(task.ToDescription())
		task.Start()
	}
}

func holdAllEnable(execArr []*core.RuleExecConfig) {
	if len(execArr) == 0 {
		return
	}

	var err error
	var procInfos []*psutil.ProcCmdInfo
	var getProcTime time.Time
	pidArr := make([]int, len(execArr), len(execArr))
	for i := 0; i < len(pidArr); i++ {
		pidArr[i] = -1
	}

	for {
		procInfos = nil

		for idx, ruleExec := range execArr {
			if pidArr[idx] > 0 && psutil.CheckPidExist(pidArr[idx]) {
				continue
			}

			if procInfos == nil {
				procInfos, err = psutil.FilterGetProcCmdInfos()
				if err != nil {
					logUtil.LoggerCommon.Error("get process error", zap.Error(err))
					time.Sleep(time.Second * 5)
					continue
				}
				getProcTime = time.Now()
			}

			pid := ruleExec.CheckAndDo(getProcTime, procInfos)
			if pid > 0 {
				pidArr[idx] = pid
			}
		}
		time.Sleep(time.Second * 2)
	}
}
