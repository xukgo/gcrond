package core

import (
	"github.com/xukgo/gcrond/psutil"
	"strings"
)

func GetProcess(procInfos []*psutil.ProcCmdInfo, execPath string, contains []string, nots []string) []*psutil.ProcCmdInfo {
	if len(contains) == 0 && len(execPath) == 0 {
		return nil
	}

	var arr = make([]*psutil.ProcCmdInfo, 0, 3)
	for _, info := range procInfos {
		cmdline := info.Cmdline
		if len(cmdline) == 0 {
			continue
		}
		exe := info.Exe
		if len(exe) == 0 {
			continue
		}

		if len(execPath) > 0 && strings.Index(exe, execPath) < 0 {
			continue
		}
		if !checkContainsAll(cmdline, contains) {
			continue
		}
		if !checkNotContainsAll(cmdline, nots) {
			continue
		}
		if !checkNotPrefixAll(cmdline, defaultExcludeCommandPrefixs) {
			continue
		}

		arr = append(arr, info)
	}
	return arr
}

//
//func GetProcessCmdLines(procInfos []*psutil.ProcCmdInfo, execPath string, contains []string, nots []string) []string {
//	if len(contains) == 0 && len(execPath) == 0 {
//		return nil
//	}
//
//	var arr = make([]string, 0, 3)
//	//infos, err := process.Processes()
//	//if err != nil {
//	//	LoggerCommon.Error("get process error", zap.Error(err))
//	//	return nil
//	//}
//	for _, info := range procInfos {
//		cmdline := info.Cmdline
//		if len(cmdline) == 0 {
//			continue
//		}
//		exe := info.Exe
//		if len(exe) == 0 {
//			continue
//		}
//		if len(execPath) > 0 && strings.Index(exe, execPath) < 0 {
//			continue
//		}
//		if !checkContainsAll(cmdline, contains) {
//			continue
//		}
//		if !checkNotContainsAll(cmdline, nots) {
//			continue
//		}
//		if !checkNotPrefixAll(cmdline, defaultExcludeCommandPrefixs) {
//			continue
//		}
//
//		arr = append(arr, cmdline)
//	}
//	return arr
//}

func checkContainsAll(str string, contains []string) bool {
	for _, containKey := range contains {
		if !strings.Contains(str, containKey) {
			return false
		}
	}
	return true
}
func checkNotContainsAll(str string, contains []string) bool {
	for _, containKey := range contains {
		if strings.Contains(str, containKey) {
			return false
		}
	}
	return true
}

func checkNotPrefixAll(str string, prefixs []string) bool {
	for _, prefix := range prefixs {
		if strings.HasPrefix(str, prefix) {
			return false
		}
	}
	return true
}
