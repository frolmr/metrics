package logger

//nolint:goimports // No clue why linter is diappointed
import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

func TestNewLogger(t *testing.T) {
	t.Run("successful logger creation", func(t *testing.T) {
		logger, err := NewLogger()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if logger == nil {
			t.Error("Expected logger to be initialized, got nil")
		}
	})

	t.Run("verify logger can log messages", func(t *testing.T) {
		core, recorded := observer.New(zapcore.DebugLevel)
		observedLogger := zap.New(core).Sugar()

		testLogger := &Logger{SugaredLogger: *observedLogger}

		testMessage := "test log message"
		testLogger.SugaredLogger.Info(testMessage)

		logs := recorded.All()
		if len(logs) != 1 {
			t.Fatalf("Expected 1 log message, got %d", len(logs))
		}

		if logs[0].Message != testMessage {
			t.Errorf("Expected message '%s', got '%s'", testMessage, logs[0].Message)
		}
	})
}

func TestLoggerMethods(t *testing.T) {
	core, recorded := observer.New(zapcore.DebugLevel)
	observedLogger := zap.New(core).Sugar()
	testLogger := &Logger{SugaredLogger: *observedLogger}

	tests := []struct {
		name     string
		logFunc  func()
		expected zapcore.Level
		message  string
	}{
		{
			name: "Debug",
			logFunc: func() {
				testLogger.SugaredLogger.Debug("debug message")
			},
			expected: zapcore.DebugLevel,
			message:  "debug message",
		},
		{
			name: "Info",
			logFunc: func() {
				testLogger.SugaredLogger.Info("info message")
			},
			expected: zapcore.InfoLevel,
			message:  "info message",
		},
		{
			name: "Warn",
			logFunc: func() {
				testLogger.SugaredLogger.Warn("warn message")
			},
			expected: zapcore.WarnLevel,
			message:  "warn message",
		},
		{
			name: "Error",
			logFunc: func() {
				testLogger.SugaredLogger.Error("error message")
			},
			expected: zapcore.ErrorLevel,
			message:  "error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logFunc()

			logs := recorded.All()
			if len(logs) < 1 {
				t.Fatal("No log messages recorded")
			}

			lastLog := logs[len(logs)-1]
			if lastLog.Level != tt.expected {
				t.Errorf("Expected level %v, got %v", tt.expected, lastLog.Level)
			}
			if lastLog.Message != tt.message {
				t.Errorf("Expected message '%s', got '%s'", tt.message, lastLog.Message)
			}
		})
	}
}
