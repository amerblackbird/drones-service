package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// go test ./internal/adapters/logger/ -v

// TestNewSimpleLogger tests creating a logger with a provided zap.Logger
func TestNewSimpleLogger(t *testing.T) {
	zapLogger, _ := zap.NewDevelopment()
	defer zapLogger.Sync()

	logger := NewSimpleLogger(zapLogger)

	if logger == nil {
		t.Fatal("Expected non-nil logger")
	}

	// Verify it's the correct type
	if _, ok := logger.(*LogService); !ok {
		t.Error("Expected logger to be of type *LogService")
	}
}

// TestNewProductionZapLogger tests creating a production logger
func TestNewProductionZapLogger(t *testing.T) {
	logger, err := NewProductionZapLogger()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if logger == nil {
		t.Fatal("Expected non-nil logger")
	}

	// Verify it's the correct type
	logService, ok := logger.(*LogService)
	if !ok {
		t.Error("Expected logger to be of type *LogService")
	}

	if logService.logger == nil {
		t.Error("Expected underlying zap.Logger to be initialized")
	}
}

// TestNewDevelopmentZapLogger tests creating a development logger
func TestNewDevelopmentZapLogger(t *testing.T) {
	logger, err := NewDevelopmentZapLogger()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if logger == nil {
		t.Fatal("Expected non-nil logger")
	}

	// Verify it's the correct type
	logService, ok := logger.(*LogService)
	if !ok {
		t.Error("Expected logger to be of type *LogService")
	}

	if logService.logger == nil {
		t.Error("Expected underlying zap.Logger to be initialized")
	}
}

// TestInfo tests the Info logging method
func TestInfo(t *testing.T) {
	// Create an observed logger to capture logs
	core, observed := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(core)
	logger := NewSimpleLogger(zapLogger)

	tests := []struct {
		name   string
		msg    string
		fields []interface{}
		want   string
	}{
		{
			name:   "info without fields",
			msg:    "test info message",
			fields: nil,
			want:   "test info message",
		},
		{
			name:   "info with fields",
			msg:    "test info with fields",
			fields: []interface{}{"key", "value", "user_id", "123"},
			want:   "test info with fields",
		},
		{
			name:   "info with single field",
			msg:    "single field test",
			fields: []interface{}{"status", "success"},
			want:   "single field test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear previous logs
			observed.TakeAll()

			// Log the message
			logger.Info(tt.msg, tt.fields...)

			// Verify log was captured
			logs := observed.All()
			if len(logs) != 1 {
				t.Fatalf("Expected 1 log entry, got %d", len(logs))
			}

			entry := logs[0]
			if entry.Message != tt.want {
				t.Errorf("Expected message %q, got %q", tt.want, entry.Message)
			}

			if entry.Level != zapcore.InfoLevel {
				t.Errorf("Expected Info level, got %v", entry.Level)
			}

			// Verify fields if provided
			if len(tt.fields) > 0 {
				if len(entry.Context) == 0 {
					t.Error("Expected fields to be logged")
				}
			}
		})
	}
}

// TestError tests the Error logging method
func TestError(t *testing.T) {
	core, observed := observer.New(zapcore.ErrorLevel)
	zapLogger := zap.New(core)
	logger := NewSimpleLogger(zapLogger)

	tests := []struct {
		name   string
		msg    string
		fields []interface{}
		want   string
	}{
		{
			name:   "error without fields",
			msg:    "test error message",
			fields: nil,
			want:   "test error message",
		},
		{
			name:   "error with fields",
			msg:    "database connection failed",
			fields: []interface{}{"error", "connection timeout", "retry_count", 3},
			want:   "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observed.TakeAll()

			logger.Error(tt.msg, tt.fields...)

			logs := observed.All()
			if len(logs) != 1 {
				t.Fatalf("Expected 1 log entry, got %d", len(logs))
			}

			entry := logs[0]
			if entry.Message != tt.want {
				t.Errorf("Expected message %q, got %q", tt.want, entry.Message)
			}

			if entry.Level != zapcore.ErrorLevel {
				t.Errorf("Expected Error level, got %v", entry.Level)
			}
		})
	}
}

// TestDebug tests the Debug logging method
func TestDebug(t *testing.T) {
	core, observed := observer.New(zapcore.DebugLevel)
	zapLogger := zap.New(core)
	logger := NewSimpleLogger(zapLogger)

	tests := []struct {
		name   string
		msg    string
		fields []interface{}
		want   string
	}{
		{
			name:   "debug without fields",
			msg:    "debug message",
			fields: nil,
			want:   "debug message",
		},
		{
			name:   "debug with fields",
			msg:    "processing request",
			fields: []interface{}{"request_id", "abc-123", "duration_ms", 45},
			want:   "processing request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observed.TakeAll()

			logger.Debug(tt.msg, tt.fields...)

			logs := observed.All()
			if len(logs) != 1 {
				t.Fatalf("Expected 1 log entry, got %d", len(logs))
			}

			entry := logs[0]
			if entry.Message != tt.want {
				t.Errorf("Expected message %q, got %q", tt.want, entry.Message)
			}

			if entry.Level != zapcore.DebugLevel {
				t.Errorf("Expected Debug level, got %v", entry.Level)
			}
		})
	}
}

// TestWarn tests the Warn logging method
func TestWarn(t *testing.T) {
	core, observed := observer.New(zapcore.WarnLevel)
	zapLogger := zap.New(core)
	logger := NewSimpleLogger(zapLogger)

	tests := []struct {
		name   string
		msg    string
		fields []interface{}
		want   string
	}{
		{
			name:   "warn without fields",
			msg:    "warning message",
			fields: nil,
			want:   "warning message",
		},
		{
			name:   "warn with fields",
			msg:    "cache miss",
			fields: []interface{}{"cache_key", "user:123", "action", "fetch_from_db"},
			want:   "cache miss",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observed.TakeAll()

			logger.Warn(tt.msg, tt.fields...)

			logs := observed.All()
			if len(logs) != 1 {
				t.Fatalf("Expected 1 log entry, got %d", len(logs))
			}

			entry := logs[0]
			if entry.Message != tt.want {
				t.Errorf("Expected message %q, got %q", tt.want, entry.Message)
			}

			if entry.Level != zapcore.WarnLevel {
				t.Errorf("Expected Warn level, got %v", entry.Level)
			}
		})
	}
}

// TestGetLogService tests retrieving the underlying zap.Logger
func TestGetLogService(t *testing.T) {
	zapLogger, _ := zap.NewDevelopment()
	defer zapLogger.Sync()

	logService := &LogService{logger: zapLogger}

	retrieved := logService.GetLogService()

	if retrieved != zapLogger {
		t.Error("Expected to retrieve the same zap.Logger instance")
	}

	if retrieved == nil {
		t.Error("Expected non-nil zap.Logger")
	}
}

// TestAsLogger tests casting Logger interface to ZapLogger
func TestAsLogger(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() interface{}
		wantOk    bool
	}{
		{
			name: "valid LogService cast",
			setupFunc: func() interface{} {
				zapLogger, _ := zap.NewDevelopment()
				return NewSimpleLogger(zapLogger)
			},
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := tt.setupFunc()

			// Cast to ports.Logger
			portsLogger, ok := logger.(interface {
				Info(msg string, fields ...interface{})
				Error(msg string, fields ...interface{})
				Debug(msg string, fields ...interface{})
				Warn(msg string, fields ...interface{})
			})

			if !ok {
				t.Fatal("Failed to cast to ports.Logger interface")
			}

			zapLogger, ok := AsLogger(portsLogger.(*LogService))

			if ok != tt.wantOk {
				t.Errorf("AsLogger() ok = %v, want %v", ok, tt.wantOk)
			}

			if tt.wantOk && zapLogger == nil {
				t.Error("Expected non-nil zap.Logger when cast is successful")
			}
		})
	}
}

// TestLogLevels tests that different log levels work correctly
func TestLogLevels(t *testing.T) {
	core, observed := observer.New(zapcore.DebugLevel)
	zapLogger := zap.New(core)
	logger := NewSimpleLogger(zapLogger)

	// Log messages at different levels
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	logs := observed.All()
	if len(logs) != 4 {
		t.Fatalf("Expected 4 log entries, got %d", len(logs))
	}

	expectedLevels := []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
	}

	for i, level := range expectedLevels {
		if logs[i].Level != level {
			t.Errorf("Log entry %d: expected level %v, got %v", i, level, logs[i].Level)
		}
	}
}

// TestLogWithComplexFields tests logging with various field types
func TestLogWithComplexFields(t *testing.T) {
	core, observed := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(core)
	logger := NewSimpleLogger(zapLogger)

	// Log with various field types
	logger.Info("complex fields test",
		"string_field", "value",
		"int_field", 42,
		"bool_field", true,
		"float_field", 3.14,
	)

	logs := observed.All()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	entry := logs[0]
	if entry.Message != "complex fields test" {
		t.Errorf("Expected message 'complex fields test', got %q", entry.Message)
	}

	// Verify that fields were logged
	if len(entry.Context) == 0 {
		t.Error("Expected context fields to be logged")
	}
}

// BenchmarkInfo benchmarks the Info logging method
func BenchmarkInfo(b *testing.B) {
	zapLogger, _ := zap.NewDevelopment()
	defer zapLogger.Sync()
	logger := NewSimpleLogger(zapLogger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark test message")
	}
}

// BenchmarkInfoWithFields benchmarks Info logging with fields
func BenchmarkInfoWithFields(b *testing.B) {
	zapLogger, _ := zap.NewDevelopment()
	defer zapLogger.Sync()
	logger := NewSimpleLogger(zapLogger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark test", "key1", "value1", "key2", 123, "key3", true)
	}
}

// BenchmarkError benchmarks the Error logging method
func BenchmarkError(b *testing.B) {
	zapLogger, _ := zap.NewDevelopment()
	defer zapLogger.Sync()
	logger := NewSimpleLogger(zapLogger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error("benchmark error message")
	}
}

// BenchmarkErrorWithFields benchmarks Error logging with fields
func BenchmarkErrorWithFields(b *testing.B) {
	zapLogger, _ := zap.NewDevelopment()
	defer zapLogger.Sync()
	logger := NewSimpleLogger(zapLogger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error("benchmark error", "error", "test error", "code", 500)
	}
}
