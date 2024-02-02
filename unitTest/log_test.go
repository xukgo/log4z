package unitTest

import (
	"fmt"
	"github.com/xukgo/log4z"
	"go.uber.org/zap"
	"log"
	"testing"
	"time"
)

// "timestamp": "2018-05-08T08:20:40.644+08:00",
func TestCallInit(t *testing.T) {
	configPath := "./conf/log4z.xml"
	loggerDict := log4z.InitLogger(configPath,
		log4z.WithTimeKey("timestamp"), log4z.WithTimeFormat("2006-01-02T15:04:05.000Z07:00"), log4z.WithCallerSkip(1),
		log4z.WithCompress(true)) //log4z.WithTimeFormat("2006-01-02T15:04:05.999Z07:00"),
	if len(loggerDict) == 0 {
		t.Fail()
	}
	logCommon, ok := loggerDict["Common"]
	if !ok {
		t.Fail()
	}

	//stdlog redirect
	//zap.RedirectStdLogAt(logCommon, zap.InfoLevel)
	zap.RedirectStdLogAt(logCommon, zap.ErrorLevel)
	//zap.RedirectStdLogAt(logCommon, zap.PanicLevel)
	//logWechat, ok := loggerDict["Wechat"]
	//if !ok {
	//	t.Fail()
	//}

	go func() {
		time.Sleep(time.Second * 3)
		log4z.SetMixConsoleLogEnable(false)
		time.Sleep(time.Second * 1)
		//os.Exit(1)
	}()

	defer func() {
		if r := recover(); r != nil {
			logCommon.Panic("panic", zap.Any("recover", r))
		}
	}()

	for {
		fmt.Println("test stdout 1")
		log.Println("test stderr 1")
		logCommon.Info("test for common appender lv Info", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
		logCommon.Warn("test for common appender lv Warn", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
		//logCommon.Error("test for common appender lv Error", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
		fmt.Println("test stdout 2")
		log.Println("test stderr 2")
		a := make([]int, 0)
		_ = a[10]
		time.Sleep(time.Second * 10)
	}
}

func TestConsoleLogger(t *testing.T) {
	logCommon := log4z.GetConsoleLogger(log4z.WithCallerSkip(0), log4z.WithMinLevel(0))
	logCommon.Debug("test for common appender lv Debug", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
	logCommon.Info("test for common appender lv Info", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
	logCommon.Warn("test for common appender lv Warn", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
	logCommon.Error("test for common appender lv Error", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))

	var gson = `{
  "imsi": "450050000000001",
  "apn": "VzWEntx1",
  "seidL": "0x10000001",
  "seidR": "0x1",
  "sxPeer": [
    "172.28.139.66"
  ],
  "ipAddrs": [
    "100.69.0.1"
  ],
  "startTime": "Thu Feb  1 10:10:32 UTC 2024",
  "throughULRate": 214689165,
  "throughDLRate": -1,
  "pdnType": "IPv4",
  "upfInstance": "vcp-tenant-ddd",
  "ppeInstance": "0",
  "upfSessLoadMbr": "0",
  "workerId": "1",
  "teidUl": "0x10000010",
  "teidDl": "0x10000110"
}`
	_ = gson
	logCommon.Info("test for json string", log4z.RawJSONField("json", []byte(gson)))
}

/*
var LoggerCommon *zap.LoggerOption //in code set the instance at a static variable
var LoggerWechat *zap.LoggerOption //in code set the instance at a static variable
func ExampleInit() {

	configPath := "./conf/log4z.xml"
	loggerMap := log4z.InitLogger(configPath)
	LoggerCommon = getLoggerOrConsole(loggerMap, "Common")
	LoggerWechat = getLoggerOrConsole(loggerMap, "Wechat")

	fmt.Println("LoggerCommon", LoggerCommon)
	fmt.Println("LoggerWechat", LoggerWechat)
}

func getLoggerOrConsole(dict map[string]*zap.LoggerOption, key string) *zap.LoggerOption {
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
*/
