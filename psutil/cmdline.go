package psutil

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Get cmdline from /proc/(pid)/cmdline
func fillFromCmdline(pid int) (string, error) {
	cmdPath := fmt.Sprintf("/proc/%d/cmdline", pid)
	cmdline, err := ioutil.ReadFile(cmdPath)
	if err != nil {
		return "", err
	}
	ret := strings.FieldsFunc(string(cmdline), func(r rune) bool {
		if r == '\u0000' {
			return true
		}
		return false
	})

	return strings.Join(ret, " "), nil
}
