package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

type LogLevel uint8

const (
	LogLevelError LogLevel = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
	LogLevelTrace
	LogLevelNone
)

type Logger struct {
	core    zerolog.Logger
	level   LogLevel
	context string
}

var ignoredLogContexts = map[string]struct{}{
	"InstanceLoader": {},
	"RoutesResolver": {},
	"RouterExplorer": {},
}

func NewExodusLogger(writer io.Writer, level string) *Logger {
	if writer == nil {
		writer = os.Stderr
	}
	resolved := parseExodusLogLevel(level)

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	zerolog.TimestampFieldName = "time"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.ErrorFieldName = "error"

	output := writer
	if !strings.EqualFold(strings.TrimSpace(os.Getenv("LOG_FORMAT")), "json") {
		output = &exodusConsoleWriter{out: writer, color: shouldColorizeLogs()}
	}

	core := zerolog.New(output).
		Level(toZerologLevel(resolved)).
		With().
		Timestamp().
		Logger()

	return &Logger{core: core, level: resolved}
}

func ResolveExodusLogLevel(configured string) string {
	configured = strings.ToLower(strings.TrimSpace(configured))
	if configured == "" {
		configured = strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL")))
	}
	if configured == "" {
		configured = strings.ToLower(strings.TrimSpace(os.Getenv("EXODUS_LOG_LEVEL")))
	}

	switch configured {
	case "trace", "verbose":
		return "trace"
	case "debug":
		return "debug"
	case "info", "log":
		return "info"
	case "warn", "warning":
		return "warn"
	case "error":
		return "error"
	case "none", "silent":
		return "none"
	}

	if parseBoolEnv(os.Getenv("ENABLE_DEBUG_LOGS")) {
		return "debug"
	}

	return "info"
}

func IsDevelopment() bool {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("NODE_ENV")))
	if env == "" {
		env = strings.ToLower(strings.TrimSpace(os.Getenv("ENV")))
	}
	return env == "development" || env == "dev"
}

func parseExodusLogLevel(level string) LogLevel {
	trimmed := strings.ToLower(strings.TrimSpace(level))
	if trimmed == "" {
		trimmed = strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL")))
	}
	if trimmed == "" {
		trimmed = strings.ToLower(strings.TrimSpace(os.Getenv("EXODUS_LOG_LEVEL")))
	}

	switch trimmed {
	case "trace", "verbose":
		return LogLevelTrace
	case "debug":
		return LogLevelDebug
	case "info", "log":
		return LogLevelInfo
	case "warn", "warning":
		return LogLevelWarn
	case "error":
		return LogLevelError
	case "none", "silent":
		return LogLevelNone
	}

	if parseBoolEnv(os.Getenv("ENABLE_DEBUG_LOGS")) {
		return LogLevelDebug
	}

	return LogLevelInfo
}

func parseBoolEnv(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "t", "true", "yes", "y", "on", "enabled":
		return true
	default:
		return false
	}
}

func shouldColorizeLogs() bool {
	if strings.TrimSpace(os.Getenv("NO_COLOR")) != "" {
		return false
	}
	if strings.EqualFold(strings.TrimSpace(os.Getenv("LOG_COLORS")), "false") {
		return false
	}
	return true
}

func (l *Logger) WithContext(context string) *Logger {
	if l == nil {
		return nil
	}
	clone := *l
	clone.context = strings.TrimSpace(context)
	return &clone
}

func (l *Logger) Enabled(level LogLevel) bool {
	return l != nil && l.level != LogLevelNone && level <= l.level
}

func (l *Logger) Log(message string, args ...any)   { l.write(LogLevelInfo, message, args...) }
func (l *Logger) Info(message string, args ...any)  { l.write(LogLevelInfo, message, args...) }
func (l *Logger) Warn(message string, args ...any)  { l.write(LogLevelWarn, message, args...) }
func (l *Logger) Error(message string, args ...any) { l.write(LogLevelError, message, args...) }
func (l *Logger) Debug(message string, args ...any) { l.write(LogLevelDebug, message, args...) }
func (l *Logger) Trace(message string, args ...any) { l.write(LogLevelTrace, message, args...) }

func (l *Logger) Fatal(message string, args ...any) {
	l.Error(message, args...)
	os.Exit(1)
}

func (l *Logger) Panic(message string, args ...any) {
	l.Error(message, args...)
	panic(message)
}

func (l *Logger) write(level LogLevel, message string, args ...any) {
	if !l.Enabled(level) {
		return
	}
	if _, ignored := ignoredLogContexts[l.context]; ignored {
		return
	}

	event := l.core.WithLevel(toZerologLevel(level))
	if l.context != "" {
		event = event.Str("context", l.context)
	}
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			event = event.Interface("extra", args[i])
			break
		}
		key, ok := args[i].(string)
		if !ok || strings.TrimSpace(key) == "" {
			event = event.Interface("extra", args[i])
			i--
			continue
		}
		trimmedKey := strings.TrimSpace(key)
		if trimmedKey == "error" {
			if e, ok := args[i+1].(error); ok {
				event = event.Err(e)
				continue
			}
		}
		event = event.Interface(trimmedKey, args[i+1])
	}
	event.Msg(message)
}

func toZerologLevel(level LogLevel) zerolog.Level {
	switch level {
	case LogLevelError:
		return zerolog.ErrorLevel
	case LogLevelWarn:
		return zerolog.WarnLevel
	case LogLevelDebug:
		return zerolog.DebugLevel
	case LogLevelTrace:
		return zerolog.TraceLevel
	case LogLevelNone:
		return zerolog.Disabled
	default:
		return zerolog.InfoLevel
	}
}

type exodusConsoleWriter struct {
	mu    sync.Mutex
	out   io.Writer
	color bool
}

func (w *exodusConsoleWriter) Write(p []byte) (int, error) {
	var event struct {
		Level   string         `json:"level"`
		Time    string         `json:"time"`
		Context string         `json:"context"`
		Message string         `json:"message"`
		Fields  map[string]any `json:"-"`
	}
	fields := map[string]any{}
	decoder := json.NewDecoder(bytes.NewReader(p))
	decoder.UseNumber()
	if err := decoder.Decode(&fields); err != nil {
		return w.writeRaw(strings.TrimRight(string(p), "\n"))
	}

	event.Level, _ = fields["level"].(string)
	event.Time, _ = fields["time"].(string)
	event.Context, _ = fields["context"].(string)
	event.Message, _ = fields["message"].(string)
	delete(fields, "level")
	delete(fields, "time")
	delete(fields, "context")
	delete(fields, "message")

	msg := strings.Trim(event.Message, "\n")
	if isBoxMessage(msg) {
		if w.color {
			msg = colorizeMultiline(msg, ansiGreen)
		}
		return w.writeRaw(msg)
	}

	line := event.Time
	level := strings.ToUpper(event.Level)
	if level == "" {
		level = "INFO"
	}
	if w.color {
		level = colorizeLevel(level)
	}
	line += " " + level
	if event.Context != "" {
		ctx := "[" + event.Context + "]"
		if w.color {
			ctx = colorize(ctx, ansiYellow)
		}
		line += " " + ctx
	}
	line += " " + event.Message
	extra := formatExtraFields(fields)
	if extra != "" {
		if w.color {
			extra = colorize(extra, ansiGray)
		}
		line += extra
	}
	return w.writeRaw(line)
}

func (w *exodusConsoleWriter) writeRaw(line string) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	_, err := fmt.Fprintln(w.out, line)
	return len(line), err
}

func isBoxMessage(message string) bool {
	return strings.HasPrefix(message, "╭") || strings.HasPrefix(message, "╔")
}

func formatExtraFields(fields map[string]any) string {
	if len(fields) == 0 {
		return ""
	}
	parts := make([]string, 0, len(fields))
	for key, value := range fields {
		parts = append(parts, key+"="+formatFieldValue(value))
	}
	return " " + strings.Join(parts, " ")
}

func formatFieldValue(value any) string {
	switch typed := value.(type) {
	case nil:
		return "null"
	case string:
		if strings.ContainsAny(typed, " \t\n\r") {
			return strconv.Quote(typed)
		}
		return typed
	case error:
		errStr := typed.Error()
		if strings.ContainsAny(errStr, " \t\n\r") {
			return strconv.Quote(errStr)
		}
		return errStr
	case fmt.Stringer:
		str := typed.String()
		if strings.ContainsAny(str, " \t\n\r") {
			return strconv.Quote(str)
		}
		return str
	default:
		data, err := json.Marshal(typed)
		if err == nil {
			s := string(data)
			if strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) {
				unquoted, unquoteErr := strconv.Unquote(s)
				if unquoteErr == nil && !strings.ContainsAny(unquoted, " \t\n\r") {
					return unquoted
				}
			}
			return s
		}
		return fmt.Sprint(typed)
	}
}

const (
	ansiReset   = "\033[0m"
	ansiRed     = "\033[31m"
	ansiYellow  = "\033[33m"
	ansiGreen   = "\033[32m"
	ansiMagenta = "\033[35m"
	ansiCyan    = "\033[36m"
	ansiGray    = "\033[90m"
)

func colorizeLevel(level string) string {
	switch level {
	case "ERROR":
		return colorize(level, ansiRed)
	case "WARN":
		return colorize(level, ansiYellow)
	case "DEBUG":
		return colorize(level, ansiMagenta)
	case "TRACE":
		return colorize(level, ansiCyan)
	default:
		return colorize(level, ansiGreen)
	}
}

func colorize(value, color string) string {
	return color + value + ansiReset
}

func colorizeMultiline(value, color string) string {
	lines := strings.Split(value, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		lines[i] = colorize(line, color)
	}
	return strings.Join(lines, "\n")
}
