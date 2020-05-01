package compon

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func GetSystemUptime() (time.Duration, error) {
	content, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}

	idx := strings.Index(string(content), " ")
	if idx <= 0 {
		return 0, fmt.Errorf("/proc/uptime格式不正确")
	}
	str := content[:idx]
	sec, err := strconv.ParseFloat(string(str), 64)
	if err != nil {
		return 0, fmt.Errorf("/proc/uptime格式不正确")
	}

	return time.Millisecond * time.Duration(sec*1000), nil
}
