package unitTest

import (
	"go.uber.org/zap"
	"log4z"
	"testing"
)

func TestCallInit(t *testing.T) {
	configPath := "./conf/log4z.xml"
	err := log4z.InitConfig(configPath)
	if err != nil {
		t.Fail()
	}
	logCommon, err := log4z.InitLogger("Common")
	if err != nil {
		t.Fail()
	}
	logWechat, err := log4z.InitLogger("Wechat")
	if err != nil {
		t.Fail()
	}
	log4z.UnintConfig()

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

/*
func ExampleInit() {
	var logCommon *zap.Logger //in code set the instance at a static variable
	var logWechat *zap.Logger //in code set the instance at a static variable
	configPath := "./conf/log4z.xml"
	err := log4z.InitConfig(configPath)
	if err != nil {
		fmt.Printf("warnning: log4z.InitConfig(configPath) configPath=%s; return err=%s\r\n", configPath, err.Error())
		fmt.Println("warnning: now set all logger to default console logger")
		logCommon = log4z.GetConsoleLogger()
		logWechat = log4z.GetConsoleLogger()
	} else {
		logCommon, err = log4z.InitLogger("Common")
		if err != nil {
			fmt.Printf("warnning: log4z.InitLogger(Common) return err=%s\r\n", err.Error())
			fmt.Println("warnning: now set logger Common to default console logger")
			logCommon = log4z.GetConsoleLogger()
		} else {
			fmt.Println("init logCommon success")
		}
		logWechat, err = log4z.InitLogger("Wechat")
		if err != nil {
			fmt.Printf("warnning: log4z.InitLogger(Wechat) return err=%s\r\n", err.Error())
			fmt.Println("warnning: now set logger Wechat to default console logger")
			logWechat = log4z.GetConsoleLogger()
		} else {
			fmt.Println("init logWechat success")
		}
	}
	log4z.UnintConfig()

	fmt.Println("logcommon", logCommon)
	fmt.Println("logWechat", logWechat)
}*/
