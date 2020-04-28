package core

import (
	"fmt"
	"github.com/xukgo/gcrond/logUtil"
	"strconv"
	"strings"
)

func ParseInterval(str string) (int64, error) {
	str = strings.ReplaceAll(str, " ", "")
	unitStr := str[len(str)-1:]
	unitStr = strings.ToLower(unitStr)
	countStr := str[:len(str)-1]
	count, err := strconv.ParseFloat(countStr, 64)
	if err != nil {
		logUtil.LoggerCommon.Error("Interval 数值格式不正确")
		return 0, err
	}
	if count < 1 {
		logUtil.LoggerCommon.Error("Interval 数值范围不正确")
		return 0, fmt.Errorf("Interval 数值范围不正确")
	}

	switch unitStr {
	case "s":
		logUtil.LoggerCommon.Info(fmt.Sprintf("解析Interval=%d%s", int64(count), "秒"))
		return int64(count), nil
	case "m":
		logUtil.LoggerCommon.Info(fmt.Sprintf("解析Interval=%v%s", count, "分钟"))
		return int64(count * 60), nil
	case "h":
		logUtil.LoggerCommon.Info(fmt.Sprintf("解析Interval=%v%s", count, "小时"))
		return int64(count * 3600), nil
	case "d":
		logUtil.LoggerCommon.Info(fmt.Sprintf("解析Interval=%v%s", count, "天"))
		return int64(count * 24 * 3600), nil
	default:
		logUtil.LoggerCommon.Error("Interval 数值单位不正确，目前只允许s/m/h/d")
		return 0, fmt.Errorf("Interval 数值单位不正确，目前只允许s/m/h/d")
	}
}
