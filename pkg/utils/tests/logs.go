package tests

import (
	"sync"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewTestLogger creates a logr.Logger instance that captures log entries for testing purposes.
// The captured logs entries are written both to Zap std logger as JSON and to CapturedLogs instance
// that can be used in unit tests to validate log messages and fields.
func NewTestLogger() (logr.Logger, *CapturedLogs) {
	capturedLogs := &CapturedLogs{}
	testCore := &testLoggerCore{logs: capturedLogs, fields: make(map[string]string)}

	// Writes logs to standard out as JSON
	stdCore := zap.NewExample().Core()

	// Create a tee core to write to both test core and standard Zap core
	teeCore := zapcore.NewTee(testCore, stdCore)

	return zapr.NewLogger(zap.New(teeCore)), capturedLogs
}

// testLoggerCore implements a zapcore.Core to captures log entries for testing
type testLoggerCore struct {
	logs   *CapturedLogs
	fields map[string]string
}

func (*testLoggerCore) Enabled(zapcore.Level) bool {
	return true
}

func (c *testLoggerCore) With(fields []zapcore.Field) zapcore.Core {
	newFields := make(map[string]string, len(c.fields)+len(fields))
	for k, v := range c.fields {
		newFields[k] = v
	}

	for _, f := range fields {
		newFields[f.Key] = fieldToString(f)
	}

	return &testLoggerCore{logs: c.logs, fields: newFields}
}

func (c *testLoggerCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return checked.AddCore(entry, c)
}

func (c *testLoggerCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	allFields := make(map[string]string, len(c.fields)+len(fields))
	for k, v := range c.fields {
		allFields[k] = v
	}

	for _, f := range fields {
		allFields[f.Key] = fieldToString(f)
	}

	c.logs.add(LogEntry{
		Message: entry.Message,
		Fields:  allFields,
	})

	return nil
}

func (*testLoggerCore) Sync() error {
	return nil
}

// LogEntry represents a captured log entry
type LogEntry struct {
	Message string
	Fields  map[string]string
}

// CapturedLogs stores captured log entries
type CapturedLogs struct {
	mu      sync.RWMutex
	entries []LogEntry
}

// FilterMessage returns all log entries matching the given message
func (o *CapturedLogs) FilterMessage(message string) []LogEntry {
	o.mu.RLock()
	defer o.mu.RUnlock()

	var result []LogEntry

	for _, e := range o.entries {
		if e.Message == message {
			result = append(result, e)
		}
	}

	return result
}

// add adds a log entry to the captured logs container
func (o *CapturedLogs) add(entry LogEntry) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.entries = append(o.entries, entry)
}

// RequireLogMessage asserts that at least one log entry with the given message was captured
// and validates optional key-value field pairs on all matching entries.
// Only the string fields are supported in this implementation.
func RequireLogMessage(t *testing.T, logs *CapturedLogs, message string, fields ...string) {
	t.Helper()

	// filter log messages by the given message
	foundMessages := logs.FilterMessage(message)
	require.NotEmpty(t, foundMessages, "expected log message not found: %s", message)

	// Validate key/value pairs for all messages
	for _, logMessage := range foundMessages {
		for i := 0; i+1 < len(fields); i += 2 {
			key, expected := fields[i], fields[i+1]

			got, ok := logMessage.Fields[key]
			if !ok || got != expected {
				t.Errorf("expected field %q=%q not found (got %q)", key, expected, got)
			}
		}
	}
}

func fieldToString(f zapcore.Field) string {
	switch f.Type {
	case zapcore.StringType:
		return f.String
	case zapcore.StringerType:
		if stringer, ok := f.Interface.(interface{ String() string }); ok {
			return stringer.String()
		}

		return ""
	default:
		return ""
	}
}
