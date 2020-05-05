package main

import (
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/process"
	_ "github.com/spf13/pflag"
	"github.com/xukgo/gcrond/compon/procUnique"
	"github.com/xukgo/gcrond/core"
	"github.com/xukgo/gcrond/logUtil"
	"github.com/xukgo/gsaber/utils/fileUtil"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
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

		procs, err := process.Processes()
		if err != nil {
			fmt.Fprintf(os.Stdout, "get process error:%s\n", err.Error())
			return
		}
		procInfos := core.GetProcess(procs, "", []string{fsearchKey}, nil)
		if len(procInfos) > 0 {
			for _, info := range procInfos {
				exe, _ := info.Exe()
				cmdline, _ := info.Cmdline()
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

	filePath := fileUtil.GetAbsUrl("conf/crond.yml")
	content, err := ioutil.ReadFile(filePath)
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

	go startTimers(conf.TimerTasks)

	if len(conf.RuleExecConfig) > 0 {
		for {
			procInfos, err := process.Processes()
			if err != nil {
				logUtil.LoggerCommon.Error("get process error", zap.Error(err))
				time.Sleep(time.Second)
				continue
			}

			for _, ruleExec := range conf.RuleExecConfig {
				if !ruleExec.Enable {
					continue
				}
				ruleExec.CheckAndDo(procInfos)
			}
			time.Sleep(time.Second)
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
