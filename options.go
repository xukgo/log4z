package log4z

import (
	"fmt"
	"github.com/xukgo/log4z/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"syscall"
	"time"
)

func newOptions(options []Option) *Options {
	opts := new(Options)
	opts.Compress = DefaultCompress
	opts.CompressDelay = DefaultCompressDelay
	opts.CallerSkip = DefaultCallerSkip

	for idx := range options {
		options[idx](opts)
	}
	if len(opts.TimeKey) == 0 {
		opts.TimeKey = DefaultTimeKey
	}
	if len(opts.LevelKey) == 0 {
		opts.LevelKey = DefaultLevelKey
	}
	if len(opts.NameKey) == 0 {
		opts.NameKey = DefaultNameKey
	}
	//if len(opts.CallerKey) == 0 {
	//	opts.CallerKey = DefaultCallerKey
	//}
	if len(opts.MessageKey) == 0 {
		opts.MessageKey = DefaultMessageKey
	}
	if len(opts.StacktraceKey) == 0 {
		opts.StacktraceKey = DefaultStacktraceKey
	}
	if len(opts.TimeFormat) == 0 {
		opts.TimeFormat = DefaultTimeFormat
	}
	return opts
}

func (this *Options) getEncoderConfig(model levelDefineXmlModel) zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        this.TimeKey,
		LevelKey:       this.LevelKey,
		NameKey:        this.NameKey,
		MessageKey:     this.MessageKey,
		StacktraceKey:  this.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     this.timeEncoder,               // 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	if model.LineRecord {
		encoderConfig.CallerKey = this.CallerKey
	} else {
		encoderConfig.CallerKey = ""
	}
	return encoderConfig
}

func (this *Options) timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(this.TimeFormat))
}

func (this *Options) createLogger(appendModel *appenderXmlModel) (*zap.Logger, error) {
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

		logPath := getAbsUrl(v.LogPath)
		maxSize := v.LogSize
		hook := lumberjack.Logger{
			Filename:      logPath,            // 日志文件路径
			MaxSize:       maxSize,            // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups:    v.MaxBackup,        // 日志文件最多保存多少个备份
			MaxAge:        v.MaxDays,          // 文件最多保存多少天
			Compress:      this.Compress,      // 是否压缩
			CompressDelay: this.CompressDelay, //N秒后才进行压缩动作
			LocalTime:     true,
		}

		//fixHookLogPath(hook)

		var WriteSync zapcore.WriteSyncer
		if v.IsConsole { //控制台和文件同时输出
			breakWriter := mixConsoleSyncSingleton.Load().(*BreakWriter)
			WriteSync = zapcore.NewMultiWriteSyncer(zapcore.AddSync(breakWriter), zapcore.AddSync(&hook))
		} else { //文件输出
			WriteSync = zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook))
		}

		var Encoder zapcore.Encoder
		if v.Encoder == "console" {
			Encoder = zapcore.NewConsoleEncoder(this.getEncoderConfig(v))
		} else {
			Encoder = zapcore.NewJSONEncoder(this.getEncoderConfig(v))
		}

		co := zapcore.NewCore(Encoder, WriteSync, levelFilter)
		core = append(core, co)
	}
	Core := zapcore.NewTee(core...)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	callerSkip := zap.AddCallerSkip(this.CallerSkip)
	// 开启文件及行号
	development := zap.Development()

	logger := zap.New(Core, caller, callerSkip, development)
	return logger, nil
}

// 对于某些没有log配置的场景下，要允许log初始化有一个执行下去的条件，就初始化成这个配置，
// 这个配置会在终端打印，相当云fmt.println，并且以console格式，常应用于testcase，不用关心log需要配置初始化
// minLevel:print log min level, -1:debug;0:info;1:warn;2:error
func (this *Options) createConsoleOnlyLogger() *zap.Logger {
	var core []zapcore.Core
	minLevel := (zapcore.Level)(this.MinLevel)
	{
		levelFilter := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= minLevel && lvl <= zapcore.WarnLevel
		})

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        this.TimeKey,
			LevelKey:       this.LevelKey,
			NameKey:        this.NameKey,
			CallerKey:      "",
			MessageKey:     this.MessageKey,
			StacktraceKey:  this.StacktraceKey,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
			EncodeTime:     this.timeEncoder,               // 时间格式
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
	}
	{
		levelFilter := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= minLevel && lvl > zapcore.WarnLevel
		})

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        this.TimeKey,
			LevelKey:       this.LevelKey,
			NameKey:        this.NameKey,
			CallerKey:      "C",
			MessageKey:     this.MessageKey,
			StacktraceKey:  this.StacktraceKey,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
			EncodeTime:     this.timeEncoder,               // 时间格式
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
	}
	Core := zapcore.NewTee(core...)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	callerSkip := zap.AddCallerSkip(this.CallerSkip)
	// 开启文件及行号
	development := zap.Development()

	logger := zap.New(Core, caller, callerSkip, development)
	return logger
}

func redirectStderr(f *os.File) {
	err := syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		log.Fatalf("Failed to redirect stderr to file: %v", err)
	}
}
