package log4z

import (
	"encoding/xml"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"strings"
)

func InitLogger(filePath string, options ...Option) map[string]*zap.Logger {
	loggerMap := make(map[string]*zap.Logger)
	localRoot, err := initConfig(filePath)
	if err == nil {
		for _, item := range localRoot.Loggers {
			loggerMap[item.Name] = getLoggerWithDefault(localRoot, item.Name, options)
		}
	}
	return loggerMap
}

func getLoggerWithDefault(confModel *confXmlRoot, loggerKey string, options []Option) *zap.Logger {
	logger, err := initLoggerByName(confModel, loggerKey, options)
	if err != nil {
		fmt.Printf("warnning: log4z.initLoggerByName(%s) return err=%s\r\n", loggerKey, err.Error())
		fmt.Printf("warnning: now set logger %s to default console logger\r\n", loggerKey)
		logger = GetConsoleLogger(options...)
	} else {
		fmt.Printf("info: init logger %s success\r\n", loggerKey)
	}

	return logger
}

func initConfig(path string) (*confXmlRoot, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	xmlroot := new(confXmlRoot)
	err = xml.Unmarshal(content, xmlroot)
	if err != nil {
		return nil, err
	}
	return xmlroot, nil
}

func parseLevel(str string) zapcore.Level {
	str = strings.ToLower(str)
	switch str {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return 8
	}
}

func initLoggerByName(confModel *confXmlRoot, name string, options []Option) (*zap.Logger, error) {
	if confModel == nil {
		return nil, fmt.Errorf("you must init logger first")
	}

	loggerXmlRoot := getLoggerFromRoot(confModel, name)
	if loggerXmlRoot == nil {
		return nil, fmt.Errorf("cannot find logger by name %s", name)
	}
	appenderXmlRoot := getAppenderFromRoot(confModel, loggerXmlRoot.AppenderName)
	if appenderXmlRoot == nil {
		return nil, fmt.Errorf("cannot find appender by name %s", loggerXmlRoot.AppenderName)
	}

	opt := newOptions(options)
	logger, err := opt.createLogger(appenderXmlRoot)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func GetConsoleLogger(opts ...Option) *zap.Logger {
	opt := newOptions(opts)
	return opt.createConsoleOnlyLogger()
}

func getLoggerFromRoot(confModel *confXmlRoot, name string) *loggerXmlModel {
	if confModel == nil {
		return nil
	}

	if confModel.Loggers == nil || len(confModel.Loggers) == 0 {
		return nil
	}

	for index := range confModel.Loggers {
		if strings.ToLower(confModel.Loggers[index].Name) == strings.ToLower(name) {
			return &confModel.Loggers[index]
		}
	}

	return nil
}

func getAppenderFromRoot(confModel *confXmlRoot, name string) *appenderXmlModel {
	if confModel == nil {
		return nil
	}

	if confModel.Appenders == nil || len(confModel.Appenders) == 0 {
		return nil
	}

	for index := range confModel.Appenders {
		if strings.ToLower(confModel.Appenders[index].Name) == strings.ToLower(name) {
			return &confModel.Appenders[index]
		}
	}

	return nil
}
