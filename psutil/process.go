package psutil

import (
	"fmt"
	"github.com/xukgo/gsaber/utils/fileUtil"
	"os"
	"strconv"
)

type PidFilterFunc func(int) bool

type ProcCmdInfo struct {
	Pid     int
	Exe     string
	Cmdline string
}

func CheckPidExist(pid int) bool {
	dir := fmt.Sprintf("/proc/%d", pid)
	exist, _ := fileUtil.Exists(dir)
	return exist
}

func FilterGetProcCmdInfos() ([]*ProcCmdInfo, error) {
	arr, err := GetProcCmdInfos(pidFilterFunc)
	if err != nil {
		return nil, err
	}
	arr = FilterProcInfos(arr)
	return arr, nil
}

func GetProcCmdInfos(pidFilterFunc PidFilterFunc) ([]*ProcCmdInfo, error) {
	pids, err := readPids()
	if err != nil {
		return nil, err
	}

	if pidFilterFunc != nil {
		cnt := 0
		for _, pid := range pids {
			if pidFilterFunc(int(pid)) {
				pids[cnt] = pid
				cnt++
			}
		}
		pids = pids[:cnt]
	}

	arr := make([]*ProcCmdInfo, 0, len(pids))
	for _, pid := range pids {
		p, err := GetProcCmdInfoByPid(int(pid))
		if err != nil {
			continue
		}
		arr = append(arr, p)
	}

	return arr, nil
}

func GetProcCmdInfoByPid(pid int) (*ProcCmdInfo, error) {
	info := new(ProcCmdInfo)
	info.Pid = pid
	info.Cmdline, _ = fillFromCmdline(pid)
	info.Exe, _ = fillFromExe(pid)
	return info, nil
}

func readPids() ([]int32, error) {
	path := "/proc"
	d, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	dirNames, err := d.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	var pidArr = make([]int32, 0, len(dirNames))
	for _, dirName := range dirNames {
		if !checkStringIsNumberic(dirName) {
			continue
		}

		pid, err := strconv.ParseInt(dirName, 10, 32)
		if err != nil {
			// if not numeric name, just skip
			continue
		}
		pidArr = append(pidArr, int32(pid))
	}

	return pidArr, nil
}
