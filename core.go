package log4z

import (
	"encoding/xml"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func InitLogger(filePath string) map[string]*zap.Logger {
	loggerMap := make(map[string]*zap.Logger)
	localRoot, err := initConfig(filePath)
	if err == nil {
		for _, item := range localRoot.Loggers {
			loggerMap[item.Name] = getLoggerWithDefault(localRoot, item.Name)
		}
	}
	return loggerMap
}

func getLoggerWithDefault(confModel *confXmlRoot, loggerKey string) *zap.Logger {
	logger, err := initLoggerByName(confModel, loggerKey)
	if err != nil {
		fmt.Printf("warnning: log4z.initLoggerByName(%s) return err=%s\r\n", loggerKey, err.Error())
		fmt.Printf("warnning: now set logger %s to default console logger\r\n", loggerKey)
		logger = GetConsoleLogger()
	} else {
		fmt.Printf("info:init logger %s success\r\n", loggerKey)
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

func initLoggerByName(confModel *confXmlRoot, name string) (*zap.Logger, error) {
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

	logger, err := createLogger(appenderXmlRoot)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func GetConsoleLogger() *zap.Logger {
	return createConsoleOnlyLogger()
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func createLogger(appendModel *appenderXmlModel) (*zap.Logger, error) {
	var core []zapcore.Core
	for _, v := range appendModel.LevelDefines {
		minLevel := parseLevel(v.MinLevel)
		if minLevel > zapcore.FatalLevel {
			return nil, fmt.Errorf("input error, please recheck MinLevel(%s) in .xml file", v.MinLevel)
		}
		maxLevel := parseLevel(v.MaxLevel)
		if maxLevel > zapcore.FatalLevel {
			return nil, fmt.Errorf("input error, please recheck MaxLevel(%s) in .xml file", v.MaxLevel)
		}
		levelFilter := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= minLevel && lvl <= maxLevel
		})

		logPath := v.LogPath
		maxSize := v.LogSize
		hook := lumberjack.Logger{
			Filename:   logPath,     // 日志文件路径
			MaxSize:    maxSize,     // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: v.MaxBackup, // 日志文件最多保存多少个备份
			MaxAge:     v.MaxDays,   // 文件最多保存多少天
			Compress:   true,        // 是否压缩
			LocalTime:  true,
		}

		var WriteSync zapcore.WriteSyncer
		if v.IsConsole { //控制台和文件同时输出
			WriteSync = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
		} else { //文件输出
			WriteSync = zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook))
		}

		var Encoder zapcore.Encoder
		if v.Encoder == "console" {
			Encoder = zapcore.NewConsoleEncoder(getEncoderConfig(v))
		} else {
			Encoder = zapcore.NewJSONEncoder(getEncoderConfig(v))
		}

		co := zapcore.NewCore(Encoder, WriteSync, levelFilter)
		core = append(core, co)
	}
	Core := zapcore.NewTee(core...)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	//caller1 := zap.AddCallerSkip(1)
	// 开启文件及行号
	development := zap.Development()

	logger := zap.New(Core, caller, development)
	return logger, nil
}

//对于某些没有log配置的场景下，要允许log初始化有一个执行下去的条件，就初始化成这个配置，
// 这个配置会在终端打印，相当云fmt.println，并且以console格式，常应用于testcase，不用关心log需要配置初始化
func createConsoleOnlyLogger() *zap.Logger {
	var core []zapcore.Core
	levelFilter := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	})

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     timeEncoder,                    // 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	var WriteSync zapcore.WriteSyncer
	WriteSync = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))

	var Encoder zapcore.Encoder
	Encoder = zapcore.NewConsoleEncoder(encoderConfig)

	co := zapcore.NewCore(Encoder, WriteSync, levelFilter)
	core = append(core, co)

	Core := zapcore.NewTee(core...)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	//caller1 := zap.AddCallerSkip(1)
	// 开启文件及行号
	development := zap.Development()

	logger := zap.New(Core, caller, development)
	return logger
}
func getEncoderConfig(model levelDefineXmlModel) zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     timeEncoder,                    // 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	if model.LineRecord {
		encoderConfig.CallerKey = "line"
	} else {
		encoderConfig.CallerKey = ""
	}
	return encoderConfig
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
