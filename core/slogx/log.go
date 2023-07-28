package slogx

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"golang.org/x/exp/slog"
)

type LogLevel int

const (
	// Silent silent log level
	Silent LogLevel = iota + 1
	// Error error log level
	Error
	// Warn warn log level
	Warn
	// Info info log level
	Info
)

type Config struct {
	LogLevel LogLevel
	App      string
}

// Writer log writer interface
type Writer interface {
	Print(context.Context, string, LogLevel, ...interface{})
}

// Interface logger interface
type LogInterface interface {
	LogMode(LogLevel) LogInterface
	Info(context.Context, string, ...interface{})
	Infof(context.Context, string, ...interface{})
	// Infof(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Warnf(context.Context, string, ...interface{})
	// Warnf(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Errorf(context.Context, string, ...interface{})
	// Errorf(context.Context, string, ...interface{})
	MustSucc(context.Context, error, ...interface{})
}

const DefaultDepth = 3

type SlogLogger struct {
	Writer
	Config
	depth     int
	addSource bool
}

func (s *SlogLogger) LogMode(level LogLevel) LogInterface {
	new_logger := *s
	new_logger.LogLevel = level
	return &new_logger
}

/*
addendum extras: []interface{}{"key1", "value1", "key2", "value2"}
*/
func (s SlogLogger) atters(args ...interface{}) []interface{} {
	forms := []interface{}{}
	if s.addSource {
		var pcs [1]uintptr
		runtime.Callers(DefaultDepth, pcs[:])
		f := runtime.CallersFrames(pcs[:])
		frame, _ := f.Next()
		forms = append(forms, "source",
			fmt.Sprintf("%s:%d", frame.File, frame.Line))
	}
	forms = append(forms, "app", s.App)
	return append(forms, args...)
}

func (s SlogLogger) Info(ctx context.Context, msg string, args ...interface{}) {

	if s.LogLevel >= Info {
		s.Writer.Print(ctx, msg, Info, s.atters(args...)...)
	}
}

func (s SlogLogger) Infof(ctx context.Context, msg string, args ...interface{}) {
	if s.LogLevel >= Info {
		s.Writer.Print(ctx, fmt.Sprintf(msg, args...), Info, s.atters()...)
	}
}

func (s SlogLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if s.LogLevel >= Warn {
		s.Writer.Print(ctx, msg, Warn, s.atters(args...)...)
	}
}

func (s SlogLogger) Warnf(ctx context.Context, msg string, args ...interface{}) {
	if s.LogLevel >= Warn {
		s.Writer.Print(ctx, fmt.Sprintf(msg, args...), Warn, s.atters()...)
	}
}

func (s SlogLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if s.LogLevel >= Error {
		s.Writer.Print(ctx, msg, Error, s.atters(args...)...)
	}
}

func (s SlogLogger) Errorf(ctx context.Context, msg string, args ...interface{}) {
	if s.LogLevel >= Error {
		s.Writer.Print(ctx, fmt.Sprintf(msg, args...), Error, s.atters()...)
	}
}

func (s SlogLogger) MustSucc(ctx context.Context, err error, args ...interface{}) {
	if err == nil {
		return
	}
	s.Writer.Print(ctx, err.Error(), Error, s.atters(args...)...)
	panic(fmt.Sprintf(err.Error(), args...))
}

func NewWithWriter(writer Writer, config Config) LogInterface {
	if config.App == "" {
		config.App = "default"
	}
	return &SlogLogger{
		Config:    config,
		Writer:    writer,
		depth:     DefaultDepth,
		addSource: true,
	}
}

var (
	Default = NewWithWriter(
		newLoggerWithJson(os.Stdout),
		Config{LogLevel: Info})
)

type logWriter struct {
	Key string
	*slog.Logger
}

func newLoggerWithJson(std *os.File) Writer {
	lw := &logWriter{}
	handler := slog.NewJSONHandler(std, &slog.HandlerOptions{
		// AddSource: true,
		/*
			is echo runtime pc
			AddSource: false
		*/
		// * modify key value
		// ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		// 	if a.Key == "time" {
		// 		return slog.Attr{Key: "ts", Value: a.Value}
		// 	}
		// 	return a
		// },
	})
	lw.Logger = slog.New(handler)
	return lw
}

func (w *logWriter) Print(ctx context.Context, msg string, lev LogLevel, args ...interface{}) {
	switch lev {
	case Silent:
		w.DebugCtx(ctx, msg, args...)
	case Error:
		w.ErrorCtx(ctx, msg, args...)
	case Warn:
		w.WarnCtx(ctx, msg, args...)
	case Info:
		w.InfoCtx(ctx, msg, args...)
	}
}
