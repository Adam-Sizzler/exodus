package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type Level int

const (
	LevelError Level = iota
	LevelWarn
	LevelInfo
	LevelHTTP
	LevelVerbose
	LevelDebug
)

type Field struct {
	Key   string
	Value any
}

type Logger struct {
	mu       *sync.Mutex
	writer   io.Writer
	level    Level
	context  string
	colors   bool
	lastTime *time.Time
}

type Options struct {
	Writer     io.Writer
	Level      string
	NodeEnv    string
	DebugLogs  string
	InstanceID string
	Colors     bool
}

var defaultLogger = New(Options{})

var contextsToIgnore = map[string]struct{}{
	"InstanceLoader": {},
	"RoutesResolver": {},
	"RouterExplorer": {},
}

func Configure(opts Options) {
	defaultLogger = New(opts)
}

func WithContext(context string) *Logger {
	return defaultLogger.WithContext(context)
}

func New(opts Options) *Logger {
	writer := opts.Writer
	if writer == nil {
		writer = os.Stderr
	}

	level := ResolveLevel(opts.NodeEnv, opts.DebugLogs, opts.Level)
	colors := opts.Colors
	if os.Getenv("NO_COLOR") != "" || strings.EqualFold(os.Getenv("LOG_COLORS"), "false") {
		colors = false
	}

	return &Logger{
		mu:       &sync.Mutex{},
		writer:   writer,
		level:    level,
		colors:   colors,
		lastTime: new(time.Time),
	}
}

func ResolveLevel(nodeEnv, debugLogs, configured string) Level {
	if strings.EqualFold(strings.TrimSpace(nodeEnv), "development") {
		return LevelDebug
	}
	if strings.EqualFold(strings.TrimSpace(debugLogs), "true") {
		return LevelDebug
	}

	switch strings.ToLower(strings.TrimSpace(configured)) {
	case "error":
		return LevelError
	case "warn", "warning":
		return LevelWarn
	case "http":
		return LevelHTTP
	case "verbose":
		return LevelVerbose
	case "debug":
		return LevelDebug
	case "info", "log", "":
		return LevelHTTP
	default:
		return LevelHTTP
	}
}

func (l *Logger) WithContext(context string) *Logger {
	clone := *l
	clone.context = strings.TrimSpace(context)
	return &clone
}

func (l *Logger) Log(message string, fields ...Field) {
	l.write(LevelInfo, "LOG", message, nil, fields...)
}

func (l *Logger) Info(message string, fields ...Field) {
	l.Log(message, fields...)
}

func (l *Logger) HTTP(message string, fields ...Field) {
	l.write(LevelHTTP, "HTTP", message, nil, fields...)
}

func (l *Logger) Warn(message string, fields ...Field) {
	l.write(LevelWarn, "WARN", message, nil, fields...)
}

func (l *Logger) Error(message string, err error, fields ...Field) {
	l.write(LevelError, "ERROR", message, err, fields...)
}

func (l *Logger) Debug(message string, fields ...Field) {
	l.write(LevelDebug, "DEBUG", message, nil, fields...)
}

func (l *Logger) Verbose(message string, fields ...Field) {
	l.write(LevelVerbose, "VERBOSE", message, nil, fields...)
}

func (l *Logger) Logf(format string, args ...any) {
	l.Log(fmt.Sprintf(format, args...))
}

func (l *Logger) Warnf(format string, args ...any) {
	l.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...any) {
	l.Error(fmt.Sprintf(format, args...), nil)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Fatal(message string, err error, fields ...Field) {
	l.Error(message, err, fields...)
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.Error(fmt.Sprintf(format, args...), nil)
	os.Exit(1)
}

func (l *Logger) Enabled(level Level) bool {
	return level <= l.level
}

func (l *Logger) write(level Level, label, message string, err error, fields ...Field) {
	if !l.Enabled(level) {
		return
	}
	if _, ignore := contextsToIgnore[l.context]; ignore {
		return
	}

	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()

	ms := "+0ms"
	if !l.lastTime.IsZero() {
		delta := now.Sub(*l.lastTime)
		if delta < 0 {
			delta = 0
		}
		ms = fmt.Sprintf("+%dms", delta.Milliseconds())
	}
	*l.lastTime = now

	prefix := fmt.Sprintf("%s %5s", now.Format("2006-01-02 15:04:05.000"), label)
	if l.colors {
		prefix = colorizeLevel(prefix, level)
	}

	context := ""
	if l.context != "" {
		context = " [" + l.context + "]"
	}

	if err != nil {
		fields = append(fields, Field{Key: "error", Value: err.Error()})
	}
	if len(fields) > 0 {
		message = strings.TrimRight(message, " ") + " " + renderFields(fields)
	}

	if strings.Contains(message, "\n") {
		fmt.Fprintf(l.writer, "%s%s %s\n", prefix, context, message)
		return
	}

	fmt.Fprintf(l.writer, "%s%s %s %s\n", prefix, context, message, ms)
}

func renderFields(fields []Field) string {
	payload := make(map[string]any, len(fields))
	for _, field := range fields {
		key := strings.TrimSpace(field.Key)
		if key == "" {
			continue
		}
		payload[key] = field.Value
	}
	if len(payload) == 0 {
		return ""
	}

	keys := make([]string, 0, len(payload))
	for key := range payload {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	ordered := make([]string, 0, len(keys))
	for _, key := range keys {
		value, err := json.Marshal(payload[key])
		if err != nil {
			value = []byte(fmt.Sprintf("%q", fmt.Sprint(payload[key])))
		}
		ordered = append(ordered, fmt.Sprintf("%q:%s", key, value))
	}
	return "{" + strings.Join(ordered, ",") + "}"
}

func colorizeLevel(text string, level Level) string {
	const reset = "\033[0m"
	switch level {
	case LevelError:
		return "\033[31m" + text + reset
	case LevelWarn:
		return "\033[33m" + text + reset
	case LevelDebug, LevelVerbose:
		return "\033[35m" + text + reset
	default:
		return "\033[32m" + text + reset
	}
}

func String(key string, value string) Field { return Field{Key: key, Value: value} }
func Int(key string, value int) Field       { return Field{Key: key, Value: value} }
func Any(key string, value any) Field       { return Field{Key: key, Value: value} }
