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
	logCommon.Info("test for common appender lv Info", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
	logCommon.Warn("test for common appender lv Warn", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
	logCommon.Error("test for common appender lv Error", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
	logWechat.Info("test for wechat appender lv Info", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
	logWechat.Warn("test for wechat appender lv Warn", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
	logWechat.Error("test for wechat appender lv Error", zap.Bool("br", true), zap.Int("int", 6001), zap.String("string", "hehehe"))
}
