package psutil

import "sync"

var pidCollectionLocker = new(sync.RWMutex)
var ignorePidCollection = make([]int, 0, 16)

func FilterProcInfos(arr []*ProcCmdInfo) []*ProcCmdInfo {
	pidCollectionLocker.Lock()
	if len(ignorePidCollection) > 10000 {
		ignorePidCollection = make([]int, 0, 16)
	}

	cnt := 0
	for _, info := range arr {
		if len(info.Exe) == 0 {
			addIgnorePid(info.Pid)
			continue
		}
		arr[cnt] = info
		cnt++
	}
	pidCollectionLocker.Unlock()

	arr = arr[:cnt]
	return arr
}

func addIgnorePid(pid int) {
	exist := false
	for idx := range ignorePidCollection {
		if ignorePidCollection[idx] == pid {
			exist = true
			break
		}
	}
	if !exist {
		ignorePidCollection = append(ignorePidCollection, pid)
	}
}

func pidFilterFunc(pid int) bool {
	pidCollectionLocker.RLock()
	defer pidCollectionLocker.RUnlock()

	if pid < 300 {
		return false
	}

	for idx := range ignorePidCollection {
		if ignorePidCollection[idx] == pid {
			return false
		}
	}
	return true
}
