package psutil

import (
	"fmt"
	"os"
)

// Get exe from /proc/(pid)/exe
func fillFromExe(pid int) (string, error) {
	exePath := fmt.Sprintf("/proc/%d/exe", pid)
	exe, err := os.Readlink(exePath)
	if err != nil {
		return "", err
	}
	return exe, nil
}
