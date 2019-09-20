package unitTest

import (
	"fmt"
	"github.com/xukgo/log4z"
	"go.uber.org/zap"
	"testing"
)

func TestCallInit(t *testing.T) {
	configPath := "./conf/log4z.xml"
	loggerDict := log4z.InitLogger(configPath)
	if len(loggerDict) == 0 {
		t.Fail()
	}
	logCommon, ok := loggerDict["Common"]
	if !ok {
		t.Fail()
	}
	logWechat, ok := loggerDict["Wechat"]
	if !ok {
		t.Fail()
	}

	logCommon.Info("test for common appender lv Info", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
	logCommon.Warn("test for common appender lv Warn", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
	logCommon.Error("test for common appender lv Error", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
	logWechat.Info("test for wechat appender lv Info", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
	logWechat.Warn("test for wechat appender lv Warn", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
	logWechat.Error("test for wechat appender lv Error", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
}

func TestConsoleLogger(t *testing.T) {
	logCommon := log4z.GetConsoleLogger()

	logCommon.Info("test for common appender lv Info", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
	logCommon.Warn("test for common appender lv Warn", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
	logCommon.Error("test for common appender lv Error", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
}

func ExampleInit() {
	var LoggerCommon *zap.Logger //in code set the instance at a static variable
	var LoggerWechat *zap.Logger //in code set the instance at a static variable

	configPath := "./conf/log4z.xml"
	loggerMap := log4z.InitLogger(configPath)
	LoggerCommon = getLoggerOrConsole(loggerMap, "Common")
	LoggerWechat = getLoggerOrConsole(loggerMap, "Wechat")

	fmt.Println("LoggerCommon", LoggerCommon)
	fmt.Println("LoggerWechat", LoggerWechat)
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
