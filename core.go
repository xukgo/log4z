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

var localRoot *ConfXmlRoot

func InitConfig(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	xmlroot := new(ConfXmlRoot)
	err = xml.Unmarshal(content, xmlroot)
	if err != nil {
		return err
	}
	localRoot = xmlroot
	return nil
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

func InitLogger(name string) (*zap.Logger, error) {
	if localRoot == nil {
		return nil, fmt.Errorf("you must init logger first")
	}

	loggerXmlRoot := getLoggerFromRoot(name)
	if loggerXmlRoot == nil {
		return nil, fmt.Errorf("cannot find logger by name %s", name)
	}
	appenderXmlRoot := getAppenderFromRoot(loggerXmlRoot.AppenderName)
	if appenderXmlRoot == nil {
		return nil, fmt.Errorf("cannot find appender by name %s", loggerXmlRoot.AppenderName)
	}

	logger, err := createLogger(appenderXmlRoot)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func createLogger(appendModel *AppenderXmlModel) (*zap.Logger, error) {
	var core []zapcore.Core
	for _, v := range appendModel.LevelDefines {
		minLevel := parseLevel(v.MinLevel)
		maxLevel := parseLevel(v.MaxLevel)
		if minLevel > zapcore.FatalLevel || maxLevel > zapcore.FatalLevel {
			return nil, fmt.Errorf("input error, please recheck MinLevel or MaxLevel in .xml file")
		}
		level := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= minLevel && lvl < maxLevel
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
		//levelPriority = append(levelPriority,level)
		//hook = append(hook,h)
		//var isConsole ,isJson bool

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

		co := zapcore.NewCore(Encoder, WriteSync, level)
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

func getEncoderConfig(model LevelDefineXmlModel) zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     TimeEncoder,                    // 时间格式
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

func getLoggerFromRoot(name string) *LoggerXmlModel {
	if localRoot == nil {
		return nil
	}

	if localRoot.Loggers == nil || len(localRoot.Loggers) == 0 {
		return nil
	}

	for index := range localRoot.Loggers {
		if strings.ToLower(localRoot.Loggers[index].Name) == strings.ToLower(name) {
			return &localRoot.Loggers[index]
		}
	}

	return nil
}
func getAppenderFromRoot(name string) *AppenderXmlModel {
	if localRoot == nil {
		return nil
	}

	if localRoot.Appenders == nil || len(localRoot.Appenders) == 0 {
		return nil
	}

	for index := range localRoot.Appenders {
		if strings.ToLower(localRoot.Appenders[index].Name) == strings.ToLower(name) {
			return &localRoot.Appenders[index]
		}
	}

	return nil
}
