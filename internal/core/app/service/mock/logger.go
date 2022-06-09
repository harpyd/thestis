package mock

import (
	"strings"
	"sync"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/pkg/deepcopy"
)

type LogStore struct {
	mu   sync.RWMutex
	logs []LogEntry
}

func (s *LogStore) flush(log LogEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logs = append(s.logs, log)
}

func (s *LogStore) flushed() []LogEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]LogEntry, len(s.logs))
	copy(result, s.logs)

	return result
}

type MemoryLogger struct {
	name string
	args map[string]interface{}

	store *LogStore
}

type LogEntry struct {
	Level   string
	Name    string
	Message string
	Err     string
	Args    map[string]interface{}
}

func NewMemoryLogger() *MemoryLogger {
	return &MemoryLogger{
		args:  make(map[string]interface{}),
		store: &LogStore{},
	}
}

func (l *MemoryLogger) FlushedLogs() []LogEntry {
	return l.store.flushed()
}

func (l *MemoryLogger) With(args ...interface{}) service.Logger {
	return &MemoryLogger{
		name:  l.name,
		args:  mergeArgs(l.args, l.argsToMap(args)),
		store: l.store,
	}
}

const nameSeparator = "."

func (l *MemoryLogger) Named(name string) service.Logger {
	return &MemoryLogger{
		name:  strings.Join([]string{l.name, name}, nameSeparator),
		args:  deepcopy.StringInterfaceMap(l.args),
		store: l.store,
	}
}

func (l *MemoryLogger) Debug(msg string, args ...interface{}) {
	l.store.flush(LogEntry{
		Name:    l.name,
		Level:   "DEBUG",
		Message: msg,
		Args:    mergeArgs(l.args, l.argsToMap(args)),
	})
}

func (l *MemoryLogger) Info(msg string, args ...interface{}) {
	l.store.flush(LogEntry{
		Name:    l.name,
		Level:   "INFO",
		Message: msg,
		Args:    mergeArgs(l.args, l.argsToMap(args)),
	})
}

func (l *MemoryLogger) Warn(msg string, args ...interface{}) {
	l.store.flush(LogEntry{
		Name:    l.name,
		Level:   "WARN",
		Message: msg,
		Args:    mergeArgs(l.args, l.argsToMap(args)),
	})
}

func (l *MemoryLogger) Error(msg string, args ...interface{}) {
	l.store.flush(LogEntry{
		Name:    l.name,
		Level:   "ERROR",
		Message: msg,
		Args:    mergeArgs(l.args, l.argsToMap(args)),
	})
}

func (l *MemoryLogger) Fatal(msg string, args ...interface{}) {
	l.store.flush(LogEntry{
		Name:    l.name,
		Level:   "FATAL",
		Message: msg,
		Args:    mergeArgs(l.args, l.argsToMap(args)),
	})
}

func (l *MemoryLogger) argsToMap(args []interface{}) map[string]interface{} {
	if len(args)%2 != 0 {
		panic("odd args number")
	}

	argsMap := make(map[string]interface{}, len(args)/2)

	for i := 0; i < len(args); i += 2 {
		key, val := args[i], args[i+1]

		if keyStr, ok := key.(string); ok {
			argsMap[keyStr] = val
		} else {
			panic("key is not string")
		}
	}

	return argsMap
}

func mergeArgs(left, right map[string]interface{}) map[string]interface{} {
	result := deepcopy.StringInterfaceMap(left)

	for k, v := range deepcopy.StringInterfaceMap(right) {
		result[k] = v
	}

	return result
}
