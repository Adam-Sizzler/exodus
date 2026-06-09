package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
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
	core       zerolog.Logger
	level      Level
	context    string
	instanceID string
	writer     *remnawaveConsoleWriter
}

type Options struct {
	Writer     io.Writer
	Level      string
	NodeEnv    string
	DebugLogs  string
	InstanceID string
	Colors     bool
	JSON       bool
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

	instanceID := strings.TrimSpace(opts.InstanceID)
	if instanceID == "" {
		instanceID = "0"
	}

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	zerolog.TimestampFieldName = "time"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.ErrorFieldName = "error"

	var output io.Writer = writer
	consoleWriter := (*remnawaveConsoleWriter)(nil)
	if !opts.JSON && !strings.EqualFold(os.Getenv("LOG_FORMAT"), "json") {
		consoleWriter = &remnawaveConsoleWriter{
			out:        writer,
			instanceID: instanceID,
			colors:     colors,
			lastTime:   new(time.Time),
		}
		output = consoleWriter
	}

	core := zerolog.New(output).
		Level(toZeroLevel(level)).
		With().
		Timestamp().
		Logger()

	return &Logger{
		core:       core,
		level:      level,
		instanceID: instanceID,
		writer:     consoleWriter,
	}
}

func ResolveLevel(nodeEnv, debugLogs, configured string) Level {
	if strings.EqualFold(strings.TrimSpace(nodeEnv), "development") || parseBool(debugLogs) {
		return LevelDebug
	}

	switch strings.ToLower(strings.TrimSpace(configured)) {
	case "error":
		return LevelError
	case "warn", "warning":
		return LevelWarn
	case "http":
		return LevelHTTP
	case "verbose", "trace":
		return LevelVerbose
	case "debug":
		return LevelDebug
	case "info", "log", "":
		return LevelInfo
	default:
		return LevelInfo
	}
}

func parseBool(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "t", "true", "yes", "y", "on", "enabled":
		return true
	default:
		return false
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

func (l *Logger) Exception(message string) {
	l.write(LevelError, "ERROR", message, nil, Field{Key: "stack", Value: []string{""}}, Field{Key: "error", Value: map[string]any{}})
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

	event := l.core.WithLevel(toZeroLevel(level)).Str("nest_level", label)
	if l.context != "" {
		event = event.Str("context", l.context)
	}
	if err != nil {
		event = event.Err(err)
	}
	for _, field := range fields {
		key := strings.TrimSpace(field.Key)
		if key == "" {
			continue
		}
		event = event.Interface(key, field.Value)
	}
	event.Msg(message)
}

func toZeroLevel(level Level) zerolog.Level {
	switch level {
	case LevelError:
		return zerolog.ErrorLevel
	case LevelWarn:
		return zerolog.WarnLevel
	case LevelVerbose:
		return zerolog.TraceLevel
	case LevelDebug:
		return zerolog.DebugLevel
	default:
		return zerolog.InfoLevel
	}
}

const maxDisplayedDelta = time.Minute

type remnawaveConsoleWriter struct {
	mu         sync.Mutex
	out        io.Writer
	instanceID string
	colors     bool
	lastTime   *time.Time
}

func (w *remnawaveConsoleWriter) Write(p []byte) (int, error) {
	line := bytes.TrimSpace(p)
	if len(line) == 0 {
		return len(p), nil
	}

	var payload map[string]any
	if err := json.Unmarshal(line, &payload); err != nil {
		_, writeErr := w.out.Write(append(line, '\n'))
		return len(p), writeErr
	}

	now := time.Now()
	if rawTime, ok := payload["time"].(string); ok {
		if parsed, err := time.ParseInLocation("2006-01-02 15:04:05.000", rawTime, time.Local); err == nil {
			now = parsed
		}
	}

	label := strings.TrimSpace(asString(payload["nest_level"]))
	if label == "" {
		label = strings.ToUpper(asString(payload["level"]))
	}
	if label == "" {
		label = "LOG"
	}

	level := levelFromString(label, asString(payload["level"]))
	context := strings.TrimSpace(asString(payload["context"]))
	message := asString(payload["message"])

	delete(payload, "time")
	delete(payload, "level")
	delete(payload, "nest_level")
	delete(payload, "context")
	delete(payload, "message")

	fields := renderPayload(payload)
	if fields != "" {
		message = strings.TrimRight(message, " ") + " - " + fields
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	ms := "+0ms"
	if !w.lastTime.IsZero() {
		delta := now.Sub(*w.lastTime)
		if delta < 0 || delta > maxDisplayedDelta {
			delta = 0
		}
		ms = fmt.Sprintf("+%dms", delta.Milliseconds())
	}
	*w.lastTime = now

	prefix := fmt.Sprintf("[#%s] %s %5s", w.instanceID, now.Format("2006-01-02 15:04:05.000"), label)
	if w.colors {
		prefix = colorizeLevel(prefix, level)
	}

	contextPart := ""
	if context != "" {
		contextPart = " [" + context + "]"
	}

	_, err := fmt.Fprintf(w.out, "%s%s %s %s\n", prefix, contextPart, message, ms)
	return len(p), err
}

func asString(value any) string {
	if value == nil {
		return ""
	}
	if str, ok := value.(string); ok {
		return str
	}
	return fmt.Sprint(value)
}

func levelFromString(label, zeroLevel string) Level {
	switch strings.ToUpper(strings.TrimSpace(label)) {
	case "ERROR", "FATAL", "PANIC":
		return LevelError
	case "WARN", "WARNING":
		return LevelWarn
	case "DEBUG":
		return LevelDebug
	case "VERBOSE", "TRACE":
		return LevelVerbose
	}
	switch strings.ToLower(strings.TrimSpace(zeroLevel)) {
	case "error", "fatal", "panic":
		return LevelError
	case "warn":
		return LevelWarn
	case "debug":
		return LevelDebug
	case "trace":
		return LevelVerbose
	default:
		return LevelInfo
	}
}

func renderPayload(payload map[string]any) string {
	if len(payload) == 0 {
		return ""
	}

	keys := make([]string, 0, len(payload))
	for key := range payload {
		if strings.TrimSpace(key) != "" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	if _, ok := payload["stack"]; ok {
		reordered := []string{"stack"}
		for _, key := range keys {
			if key != "stack" {
				reordered = append(reordered, key)
			}
		}
		keys = reordered
	}

	ordered := make([]string, 0, len(keys))
	for _, key := range keys {
		ordered = append(ordered, fmt.Sprintf("%s: %s", key, renderValue(payload[key])))
	}
	return "{ " + strings.Join(ordered, ", ") + " }"
}

func renderValue(value any) string {
	switch typed := value.(type) {
	case string:
		return fmt.Sprintf("%q", typed)
	case []string:
		items := make([]string, 0, len(typed))
		for _, item := range typed {
			items = append(items, fmt.Sprintf("%q", item))
		}
		return "[ " + strings.Join(items, ", ") + " ]"
	case []any:
		items := make([]string, 0, len(typed))
		for _, item := range typed {
			items = append(items, renderValue(item))
		}
		return "[ " + strings.Join(items, ", ") + " ]"
	case map[string]any:
		if len(typed) == 0 {
			return "{}"
		}
	}

	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%q", fmt.Sprint(value))
	}
	return string(encoded)
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
