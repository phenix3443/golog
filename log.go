package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	DLogger *zap.SugaredLogger // 诊断日志
	SLogger *zap.Logger        // 统计日志
)

func NewDiagosisLogger(cfg *zap.Config) (s *zap.SugaredLogger, err error) {
	l, err := cfg.Build()
	if err != nil {
		return
	}
	s = l.Sugar()

	return
}

func NewStaticLogger(cfg *zap.Config) (l *zap.Logger, err error) {
	l, err = cfg.Build()
	return
}

func DLoggerCfg() (cfg *zap.Config) {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	atomLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg = &zap.Config{
		Level:       atomLevel,
		Development: true,
		// DisableCaller:     false,
		// DisableStacktrace: false,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		// InitialFields:     map[string]interface{}{"foo": "bar"},
	}

	return
}

func SLoggerCfg() (cfg *zap.Config) {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	atomLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg = &zap.Config{
		Level:             atomLevel,
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Encoding:          "json",
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stdout", "/tmp/stdout"},
		ErrorOutputPaths:  []string{"stdout"},
		InitialFields:     map[string]interface{}{"foo": "bar"},
	}

	return
}

func init() {
	var err error
	if DLogger, err = NewDiagosisLogger(DLoggerCfg()); err != nil {
		fmt.Printf("初始化诊断日志系统失败,err=%v\n", err)
	} else {
		DLogger.Info("诊断日志系统初始化成功")
	}

	if SLogger, err = NewStaticLogger(SLoggerCfg()); err != nil {
		fmt.Printf("初始化统计日志系统失败,err=%v\n", err)
	} else {
		DLogger.Info("统计日志系统初始化成功")
	}
}
