package log

import (
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"os"
)

var jsonLogger *Logger
var textLogger *Logger

// Logger 封装后的日志结构
type Logger struct {
	logger *logrus.Logger
}

// 初始化函数
func init() {
	logrusInit()
}

func logrusInit() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "",    // 使用默认的时间戳格式
		DisableTimestamp:  false, // 不禁用时间戳
		DisableHTMLEscape: false, // 不禁用HTML转义
		DataKey:           "",    // 不使用DataKey
		FieldMap:          nil,   // 使用默认的字段映射
		CallerPrettyfier:  nil,   // 不自定义Caller格式
		PrettyPrint:       true,  // 美化JSON输出
	})

	logger.SetOutput(os.Stdout)
	// 将封装的 logger 赋值到全局变量中
	jsonLogger = &Logger{logger: logger}

	logger2 := logrus.New()
	logger2.SetFormatter(&logrus.TextFormatter{
		ForceColors:               true,
		DisableColors:             false,
		ForceQuote:                false,
		DisableQuote:              false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           "",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		PadLevelText:              true,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier:          nil,
	})
	logger2.SetOutput(colorable.NewColorableStdout())
	textLogger = &Logger{logger: logger2}
}

// GetJsonLogger 获取封装后的 Json Logger 实例
func GetJsonLogger() *Logger {
	return jsonLogger
}

// GetTextLogger 获取封装后的文本Logger实例
func GetTextLogger() *Logger { return textLogger }

// WithFields 为 Logger 添加结构化字段并返回新的 Logger
// 简化的WithFields方法，使用变长参数
func (l *Logger) WithFields(fields ...interface{}) *Logger {
	if len(fields)%2 != 0 {
		l.logger.Warn("Invalid number of arguments for WithFields")
		return l
	}

	fieldMap := logrus.Fields{}
	for i := 0; i < len(fields); i += 2 {
		key, okKey := fields[i].(string)
		if !okKey {
			l.logger.Warn("Key must be a string")
			continue
		}
		fieldMap[key] = fields[i+1]
	}

	return &Logger{
		logger: l.logger.WithFields(fieldMap).Logger,
	}
}

// Info 级别的日志记录方法
func (l *Logger) Info(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Info(msg)
	}
	l.logger.Infof(msg, args)
}

// Warn 级别的日志记录方法
func (l *Logger) Warn(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Warn(msg)
	}
	l.logger.Warnf(msg, args)
}

// Error 级别的日志记录方法
func (l *Logger) Error(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Error(msg)
	}
	l.logger.Errorf(msg, args)
}

// Fatal 级别的日志记录方法
func (l *Logger) Fatal(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Fatal(msg)
	}
	l.logger.Fatalf(msg, args)
}
