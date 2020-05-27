package psutil

import (
	"fmt"
	"testing"
)

func TestGetProcCmdInfo(t *testing.T) {
	arr, _ := GetProcCmdInfos(pidFilterFunc)
	arr = FilterProcInfos(arr)
	fmt.Println(len(arr))
}
