package psutil

import (
	"github.com/xukgo/gcrond/compon"
	"log"
	"sync"
	"time"
)

var uptimeLocker = new(sync.Mutex)
var lastUptime time.Duration
var lastCheckTime = time.Unix(0, 0)

func GetUptime() (time.Duration, error) {
	uptimeLocker.Lock()
	defer uptimeLocker.Unlock()

	dur := time.Since(lastCheckTime)
	if dur.Hours() < 6 {
		return dur + lastUptime, nil
	}

	upduration, err := compon.GetSystemUptime()
	if err != nil {
		log.Println("get system uptime error:", err.Error())
		return 0, err
	}
	lastUptime = upduration
	lastCheckTime = time.Now()
	return upduration, nil
}
