package core

import (
	"github.com/shirou/gopsutil/process"
	"go.uber.org/zap"
	"strings"
)

/* {"procDir": "/opt/newcc", "cmdline": "./newcc", "exe": "/opt/newcc/newcc", "name": "newcc", "username": "hermes"}
cmdline是输入到命令行原文，exe是可执行文件的实际路径，穿透了软连接的
*/
func GetProcessInfos() ([]*process.Process, error) {
	infos, err := process.Processes()
	for _, info := range infos {
		cmdline, _ := info.Cmdline()
		if len(cmdline) == 0 {
			continue
		}
		procDir, _ := info.Cwd()
		exe, _ := info.Exe()
		name, _ := info.Name()
		username, _ := info.Username()

		if strings.Index(exe, "newcc") < 0 {
			continue
		}
		LoggerCommon.Info("", zap.String("procDir", procDir), zap.String("cmdline", cmdline), zap.String("exe", exe),
			zap.String("name", name), zap.String("username", username))
	}
	return infos, err
}

func GetProcessCmdLines(contains []string, nots []string) []string {
	if len(contains) == 0 {
		return nil
	}

	var arr = make([]string, 0, 3)
	infos, err := process.Processes()
	if err != nil {
		LoggerCommon.Error("get process error", zap.Error(err))
		return nil
	}
	for _, info := range infos {
		cmdline, _ := info.Cmdline()
		if len(cmdline) == 0 {
			continue
		}
		//exe,_ := info.Exe()
		if strings.HasPrefix(cmdline, "cat ") {
			continue
		}
		if strings.HasPrefix(cmdline, "zcat ") {
			continue
		}
		if strings.HasPrefix(cmdline, "vim ") {
			continue
		}
		if strings.HasPrefix(cmdline, "tail ") {
			continue
		}
		if strings.HasPrefix(cmdline, "gdb ") {
			continue
		}
		if strings.Contains(cmdline, "grep ") {
			continue
		}
		if !checkContainsAll(cmdline, contains) {
			continue
		}
		if !checkNotContainsAll(cmdline, nots) {
			continue
		}

		arr = append(arr, cmdline)
	}
	return arr
}

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