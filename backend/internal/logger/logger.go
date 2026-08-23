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
	LevelVerbose
	LevelDebug
)

type Field struct {
	Key   string
	Value any
}

type Logger struct {
	core    zerolog.Logger
	level   Level
	context string
	writer  *exodusConsoleWriter
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

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	zerolog.TimestampFieldName = "time"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.ErrorFieldName = "error"

	var output io.Writer = writer
	consoleWriter := (*exodusConsoleWriter)(nil)
	if !opts.JSON && !strings.EqualFold(os.Getenv("LOG_FORMAT"), "json") {
		consoleWriter = &exodusConsoleWriter{
			out:    writer,
			colors: colors,
		}
		output = consoleWriter
	}

	core := zerolog.New(output).
		Level(toZeroLevel(level)).
		With().
		Timestamp().
		Logger()

	return &Logger{
		core:   core,
		level:  level,
		writer: consoleWriter,
	}
}

func ResolveLevel(nodeEnv, debugLogs, configured string) Level {
	trimmed := strings.ToLower(strings.TrimSpace(configured))
	if trimmed == "" {
		trimmed = strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL")))
	}

	switch trimmed {
	case "error":
		return LevelError
	case "warn", "warning":
		return LevelWarn
	case "verbose", "trace":
		return LevelVerbose
	case "debug":
		return LevelDebug
	case "info", "log":
		return LevelInfo
	}

	if parseBool(debugLogs) {
		return LevelDebug
	}

	return LevelInfo
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
	l.write(LevelInfo, message, nil, fields...)
}

func (l *Logger) Info(message string, fields ...Field) {
	l.Log(message, fields...)
}

func (l *Logger) HTTP(message string, fields ...Field) {
	l.Debug(message, fields...)
}

func (l *Logger) Warn(message string, fields ...Field) {
	l.write(LevelWarn, message, nil, fields...)
}

func (l *Logger) Error(message string, err error, fields ...Field) {
	l.write(LevelError, message, err, fields...)
}

func (l *Logger) Debug(message string, fields ...Field) {
	l.write(LevelDebug, message, nil, fields...)
}

func (l *Logger) Verbose(message string, fields ...Field) {
	l.write(LevelVerbose, message, nil, fields...)
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

func (l *Logger) ConfigError(message string) {
	l.write(LevelError, message, nil, Field{Key: "stack", Value: []string{""}}, Field{Key: "error", Value: map[string]any{}})
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

func (l *Logger) write(level Level, message string, err error, fields ...Field) {
	if !l.Enabled(level) {
		return
	}
	if _, ignore := contextsToIgnore[l.context]; ignore {
		return
	}

	event := l.core.WithLevel(toZeroLevel(level))
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

const (
	ansiReset  = "[0m"
	ansiRed    = "[31m"
	ansiGreen  = "[32m"
	ansiYellow = "[33m"
	ansiCyan   = "[36m"
)

type exodusConsoleWriter struct {
	mu     sync.Mutex
	out    io.Writer
	colors bool
}

func (w *exodusConsoleWriter) Write(p []byte) (int, error) {
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

	label := consoleLevelLabel(asString(payload["level"]))
	level := levelFromString(label, asString(payload["level"]))
	context := strings.TrimSpace(asString(payload["context"]))
	message := asString(payload["message"])

	delete(payload, "time")
	delete(payload, "level")
	delete(payload, "context")
	delete(payload, "message")

	fields := renderPayload(payload)
	if fields != "" {
		message = strings.TrimRight(message, " ") + " " + fields
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if isBoxMessage(message) {
		if w.colors {
			message = colorizeMultiline(message, ansiGreen)
		}
		return w.writeRaw(message)
	}

	timePart := now.Format("2006-01-02 15:04:05.000")
	levelPart := label
	contextPart := ""
	if context != "" {
		contextPart = " [" + context + "]"
	}
	if w.colors {
		levelPart = colorizeLevel(label, level)
		if contextPart != "" {
			contextPart = colorize(contextPart, ansiYellow)
		}
	}

	_, err := fmt.Fprintf(w.out, "%s %s%s %s\n", timePart, levelPart, contextPart, message)
	return len(p), err
}

func (w *exodusConsoleWriter) writeRaw(message string) (int, error) {
	message = strings.Trim(message, "\n")
	if message == "" {
		return 0, nil
	}
	return fmt.Fprintln(w.out, message)
}

func colorizeMultiline(value, color string) string {
	if value == "" {
		return value
	}

	lines := strings.Split(value, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		lines[i] = colorize(line, color)
	}
	return strings.Join(lines, "\n")
}

func isBoxMessage(message string) bool {
	message = strings.TrimSpace(message)
	return strings.HasPrefix(message, "╭") || strings.HasPrefix(message, "╔")
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

func consoleLevelLabel(zeroLevel string) string {
	switch strings.ToLower(strings.TrimSpace(zeroLevel)) {
	case "error", "fatal", "panic":
		return "ERROR"
	case "warn":
		return "WARN"
	case "debug":
		return "DEBUG"
	case "trace":
		return "TRACE"
	case "info":
		return "INFO"
	default:
		return "INFO"
	}
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
	switch level {
	case LevelError:
		return colorize(text, ansiRed)
	case LevelWarn:
		return colorize(text, ansiYellow)
	case LevelDebug, LevelVerbose:
		return colorize(text, ansiCyan)
	default:
		return colorize(text, ansiGreen)
	}
}

func colorize(text, color string) string {
	return color + text + ansiReset
}

func String(key string, value string) Field { return Field{Key: key, Value: value} }
func Int(key string, value int) Field       { return Field{Key: key, Value: value} }
func Any(key string, value any) Field       { return Field{Key: key, Value: value} }
