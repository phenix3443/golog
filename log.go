package log

import (
	"fmt"
	"log"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	DLogger     *zap.SugaredLogger // 诊断日志
	DLoggerPort = 9876             // 监听端口
	SLogger     *zap.Logger        // 统计日志
	SLoggerPort = 9875             // 监听端口
)

// NewDiagosisLogger 创建新的诊断日志记录器，使用控制台记录格式
func NewDiagosisLogger(cfg *zap.Config) (s *zap.SugaredLogger, err error) {
	l, err := cfg.Build()
	if err != nil {
		return
	}
	s = l.Sugar()

	// 动态改变日志等级
	http.HandleFunc("/dlogger", cfg.Level.ServeHTTP)

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", DLoggerPort), nil); err != nil {
			log.Fatalf("诊断日志端口监听失败,err=%v", err)
		}
	}()

	return
}

// NewStaticLogger创建新的统计日志记录器，使用json格式
func NewStaticLogger(cfg *zap.Config) (l *zap.Logger, err error) {
	l, err = cfg.Build()
	// 动态改变日志等级
	http.HandleFunc("/slogger", cfg.Level.ServeHTTP)

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", SLoggerPort), nil); err != nil {
			log.Fatalf("数据日志端口监听失败,err=%v", err)
		}
	}()
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

// Init 初始化函数，调用该函数可以直接初始化log package
func Init() {
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
