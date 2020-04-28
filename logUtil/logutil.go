package logUtil

import (
	"fmt"
	"github.com/xukgo/gsaber/utils/fileUtil"
	"github.com/xukgo/log4z"
	"go.uber.org/zap"
)

var LoggerCommon *Logger
var LoggerBll *Logger

func init() {
	//confPath := fileUtil.GetAbsUrl("conf/log4z.xml")
	//loggerMap := log4z.InitLogger(confPath)
	//LoggerCommon = getLoggerOrConsole(loggerMap, "Common")
	//LoggerBll = getLoggerOrConsole(loggerMap, "Bll")
	//LoggerPush = getLoggerOrConsole(loggerMap, "Push")
	//LoggerSdr = getLoggerOrConsole(loggerMap, "Sdr")
	//LoggerWrongCdr = getLoggerOrConsole(loggerMap, "WrongCdr")

	confPath := fileUtil.GetAbsUrl("conf/log4z.xml")
	loggerMap := log4z.InitLogger(confPath,
		log4z.WithTimeKey("timestamp"), log4z.WithTimeFormat("2006-01-02 15:04:05.999"))
	elkLogger := getLoggerOrConsole(loggerMap, "Elk")

	LoggerCommon = newLogger(elkLogger, INNER_MODULE_COMMON)
	LoggerBll = newLogger(elkLogger, INNER_MODULE_BLL)
}
func getLoggerOrConsole(dict map[string]*zap.Logger, key string) *zap.Logger {
	logger, ok := dict[key]
	if ok {
		fmt.Printf("info: get logger %s success\r\n", key)
	} else {
		fmt.Printf("warnning: log4z get logger (%s) failed\r\n", key)
		fmt.Printf("warnning: now set logger %s to default console logger\r\n", key)
		logger = log4z.GetConsoleLogger()
	}
	return logger
}
