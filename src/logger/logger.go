package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"go-sqli-lab/src/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Level 日志级别
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger 日志接口
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	Close() error
}

// defaultLogger 默认日志实现
type defaultLogger struct {
	level  Level
	format string
	output io.WriteCloser
	module string
	mu     sync.Mutex
}

// New 创建日志实例
func New(cfg config.LogConfig) (Logger, error) {
	level := parseLevel(cfg.Level)

	var output io.WriteCloser
	if cfg.Output == "stdout" {
		output = os.Stdout
	} else {
		// 确保日志目录存在
		dir := filepath.Dir(cfg.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("创建日志目录失败: %w", err)
		}

		// 使用lumberjack实现滚动日志
		output = &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
	}

	return &defaultLogger{
		level:  level,
		format: cfg.Format,
		output: output,
		module: "APP",
	}, nil
}

func parseLevel(level string) Level {
	switch strings.ToLower(level) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn", "warning":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO
	}
}

func (l *defaultLogger) log(level Level, msg string, keysAndValues ...interface{}) {
	if level < l.level {
		return
	}

	// 获取调用者信息
	_, file, line, ok := runtime.Caller(2)
	function := "unknown"
	if ok {
		function = getFunctionName()
		file = filepath.Base(file)
	}

	// 获取goroutine ID
	gid := getGoroutineID()

	// 构建日志内容
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// 格式化keysAndValues
	var fields string
	if len(keysAndValues) > 0 {
		var sb strings.Builder
		for i := 0; i < len(keysAndValues); i += 2 {
			if i > 0 {
				sb.WriteString(", ")
			}
			key := fmt.Sprintf("%v", keysAndValues[i])
			var value string
			if i+1 < len(keysAndValues) {
				value = fmt.Sprintf("%v", keysAndValues[i+1])
			} else {
				value = "MISSING"
			}
			sb.WriteString(key)
			sb.WriteString("=")
			sb.WriteString(value)
		}
		fields = sb.String()
	}

	// 应用格式
	logLine := l.format
	logLine = strings.ReplaceAll(logLine, "%timestamp%", timestamp)
	logLine = strings.ReplaceAll(logLine, "%level%", level.String())
	logLine = strings.ReplaceAll(logLine, "%goroutine%", gid)
	logLine = strings.ReplaceAll(logLine, "%module%", l.module)
	logLine = strings.ReplaceAll(logLine, "%function%", function)
	logLine = strings.ReplaceAll(logLine, "%file%", file)
	logLine = strings.ReplaceAll(logLine, "%line%", fmt.Sprintf("%d", line))
	logLine = strings.ReplaceAll(logLine, "%message%", msg)

	if fields != "" {
		logLine = logLine + " | " + fields
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintln(l.output, logLine)

	if level == FATAL {
		os.Exit(1)
	}
}

func (l *defaultLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.log(DEBUG, msg, keysAndValues...)
}

func (l *defaultLogger) Info(msg string, keysAndValues ...interface{}) {
	l.log(INFO, msg, keysAndValues...)
}

func (l *defaultLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.log(WARN, msg, keysAndValues...)
}

func (l *defaultLogger) Error(msg string, keysAndValues ...interface{}) {
	l.log(ERROR, msg, keysAndValues...)
}

func (l *defaultLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.log(FATAL, msg, keysAndValues...)
}

func (l *defaultLogger) Close() error {
	if l.output != os.Stdout {
		return l.output.Close()
	}
	return nil
}

// getFunctionName 获取调用函数名
func getFunctionName() string {
	pc, _, _, ok := runtime.Caller(3)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	name := fn.Name()
	// 只返回函数名，不包含包路径
	if idx := strings.LastIndex(name, "."); idx != -1 {
		return name[idx+1:]
	}
	return name
}

// getGoroutineID 获取当前goroutine ID
func getGoroutineID() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	return idField
}
