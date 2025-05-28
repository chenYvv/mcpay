// pkg/logger/daily_writer.go
package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// 全局变量
var (
	Logger      *slog.Logger
	dailyWriter *DailyWriter
)

// DailyWriter 按天切换的日志写入器
type DailyWriter struct {
	dir    string   // 日志目录
	prefix string   // 文件名前缀
	suffix string   // 文件扩展名
	file   *os.File // 当前日志文件
	date   string   // 当前日期
	mu     sync.Mutex
}

// NewDailyWriter 创建新的按天日志写入器
func NewDailyWriter(dir, prefix, suffix string) *DailyWriter {
	os.MkdirAll(dir, 0755)

	writer := &DailyWriter{
		dir:    dir,
		prefix: prefix,
		suffix: suffix,
	}

	writer.rotate()
	return writer
}

func (w *DailyWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 检查是否需要轮转
	currentDate := time.Now().Format("2006-01-02")
	if currentDate != w.date {
		w.rotate()
	}

	return w.file.Write(p)
}

func (w *DailyWriter) rotate() {
	// 关闭旧文件
	if w.file != nil {
		w.file.Close()
	}

	// 创建新文件
	w.date = time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s%s", w.date, w.suffix)
	filePath := filepath.Join(w.dir, filename)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("无法创建日志文件: %v", err))
	}

	w.file = file
}

func (w *DailyWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

// InitSlog 初始化 slog 日志系统
func InitSlog(serviceName string, debug bool) {
	// 创建日志目录
	logDir := filepath.Join("logs", serviceName)
	dailyWriter = NewDailyWriter(logDir, serviceName, ".log")

	// 设置输出
	var writer io.Writer

	if !debug {
		writer = dailyWriter
	} else {
		writer = io.MultiWriter(os.Stdout, dailyWriter)
	}

	// 生产环境和开发环境使用不同配置
	var handler slog.Handler
	if !debug {
		// 生产环境：简洁格式
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: false, // 关键：关闭 source
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.String("time", time.Now().Format("2006-01-02 15:04:05"))
				}
				return a
			},
		})
	} else {
		// 开发环境：详细格式，但简化 source
		handler = slog.NewTextHandler(writer, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.String("time", time.Now().Format("15:04:05"))
				}
				if a.Key == slog.SourceKey {
					if source, ok := a.Value.Any().(*slog.Source); ok {
						// 只保留文件名和行号
						return slog.String("source", fmt.Sprintf("%s:%d",
							filepath.Base(source.File), source.Line))
					}
				}
				return a
			},
		})
	}

	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

// CleanupLogger 清理日志资源
func CleanupLogger() {
	if dailyWriter != nil {
		slog.Info("正在关闭日志系统...")
		dailyWriter.Close()
	}
}

// 便捷方法
func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	Logger.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	Logger.Warn(msg, args...)
}

func Fatal(msg string, args ...any) {
	Logger.Error(msg, args...)
	os.Exit(1)
}

// 兼容原有的 String 方法
func InfoString(module, name, msg string) {
	Logger.Info(module, slog.String(name, msg))
}

func ErrorString(module, name, msg string) {
	Logger.Error(module, slog.String(name, msg))
}

func DebugString(module, name, msg string) {
	Logger.Debug(module, slog.String(name, msg))
}

func WarnString(module, name, msg string) {
	Logger.Warn(module, slog.String(name, msg))
}
