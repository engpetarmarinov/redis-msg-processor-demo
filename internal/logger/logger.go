package logger

import (
	"log/slog"
	"os"
	"runtime/debug"
	"strings"
	"sync/atomic"
)

var stdLogger atomic.Pointer[Logger]
var errLogger atomic.Pointer[Logger]

// Logger wraps slog.Logger
type Logger struct {
	*slog.Logger
}

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func newLogLevel(level string) Level {
	level = strings.ToUpper(level)
	switch level {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}

func getSLogLevel(level Level) slog.Level {
	switch level {
	case DEBUG:
		return slog.LevelDebug
	case INFO:
		return slog.LevelInfo
	case WARN:
		return slog.LevelWarn
	case ERROR:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type ConfigOpt struct {
	Level Level
}

// WithLevel sets a threshold to the logger options
func (opt *ConfigOpt) WithLevel(level string) *ConfigOpt {
	opt.Level = newLogLevel(level)
	return opt
}

// NewConfigOpt creates ConfigOpt
func NewConfigOpt() *ConfigOpt {
	return &ConfigOpt{}
}

func initStdLogger(cfgOpt *ConfigOpt) {
	if cfgOpt == nil {
		cfgOpt = NewConfigOpt()
	}

	slogLevel := getSLogLevel(cfgOpt.Level)
	lvl := new(slog.LevelVar)
	lvl.Set(slogLevel)
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))

	stdLogger.Store(&Logger{l})
}

func initErrLogger() {
	l := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:       slog.LevelError,
		ReplaceAttr: replaceAttr,
	}))

	errLogger.Store(&Logger{l})
}

func replaceAttr(_ []string, a slog.Attr) slog.Attr {
	switch a.Value.Kind() {
	case slog.KindAny:
		switch v := a.Value.Any().(type) {
		case error:
			a.Value = fmtErr(v)
		}
	}

	return a
}

func fmtErr(err error) slog.Value {
	var groupValues []slog.Attr
	groupValues = append(groupValues, slog.String("msg", err.Error()))
	groupValues = append(groupValues, slog.Any("trace", debug.Stack()))

	return slog.GroupValue(groupValues...)
}

// Init sets the initial values of the logger, initializing stdout and stderr loggers
func Init(cfgOpt *ConfigOpt) {
	initStdLogger(cfgOpt)
	initErrLogger()
}

// Debug logs debug msg and args
func Debug(msg string, args ...any) {
	stdLogger.Load().Debug(msg, args...)
}

// Info logs info msg and args
func Info(msg string, args ...any) {
	stdLogger.Load().Info(msg, args...)
}

// Warn logs warning msg and args
func Warn(msg string, args ...any) {
	stdLogger.Load().Warn(msg, args...)
}

// Error logs error msg and args to the stderr fd
func Error(msg string, args ...any) {
	errLogger.Load().Error(msg, args...)
}
