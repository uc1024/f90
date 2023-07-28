package slogx

import (
	"context"
	"os"
	"testing"
)

// logger.Info(os.UserCacheDir())
// deg := struct{ Id int }{Id: int(1)}
// logger.Info("deg",deg)

type TestInfo struct {
	FIX   string
	Level string
}

func TestLogJson(t *testing.T) {
	// logger := Default.LogMode(Info)
	Default.Info(context.Background(), "test")
	Default.Info(context.Background(), "test log ", "test", TestInfo{FIX: "fix-log", Level: "info"})
	Default.Infof(context.Background(), "test xxxxlog %+v", TestInfo{FIX: "fix-log", Level: "info"})

	logger_test := NewWithWriter(
		newLoggerWithJson(os.Stdout),
		Config{LogLevel: Error})
	logger_test.Info(context.Background(), "test log ", "logger_test", TestInfo{FIX: "fix-log", Level: "info"})
	logger_test.Error(context.Background(), "test log ", "logger_test", TestInfo{FIX: "fix-log", Level: "error"})
	logger_test.Errorf(context.Background(), "test log %v", TestInfo{FIX: "fix-log", Level: "errorf"})
}
