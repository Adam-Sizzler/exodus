package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
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

type loggerCore struct {
	mu     sync.Mutex
	writer io.Writer
	level  LogLevel
	color  bool
}

type Logger struct {
	core    *loggerCore
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
	return &Logger{core: &loggerCore{
		writer: writer,
		level:  parseExodusLogLevel(level),
		color:  shouldColorizeLogs(),
	}}
}

func ResolveExodusLogLevel(configured string) string {
	if IsDevelopment() {
		return "debug"
	}

	configured = strings.ToLower(strings.TrimSpace(configured))
	if configured != "" {
		switch configured {
		case "info", "log", "warn", "warning", "error", "none", "silent":
			return configured
		}
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
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "trace", "verbose":
		return LogLevelTrace
	case "debug":
		return LogLevelDebug
	case "", "info", "log":
		return LogLevelInfo
	case "warn", "warning":
		return LogLevelWarn
	case "error":
		return LogLevelError
	case "none", "silent":
		return LogLevelNone
	default:
		if IsDevelopment() {
			return LogLevelDebug
		}
		return LogLevelInfo
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
	context = strings.TrimSpace(context)
	return &Logger{core: l.core, context: context}
}

func (l *Logger) Enabled(level LogLevel) bool {
	if l == nil || l.core == nil || l.core.level == LogLevelNone {
		return false
	}
	return level <= l.core.level
}

func (l *Logger) Log(message string, args ...any)  { l.write(LogLevelInfo, "LOG", message, args...) }
func (l *Logger) Info(message string, args ...any) { l.write(LogLevelInfo, "LOG", message, args...) }
func (l *Logger) Warn(message string, args ...any) { l.write(LogLevelWarn, "WARN", message, args...) }
func (l *Logger) Error(message string, args ...any) {
	l.write(LogLevelError, "ERROR", message, args...)
}
func (l *Logger) Debug(message string, args ...any) {
	l.write(LogLevelDebug, "DEBUG", message, args...)
}
func (l *Logger) Trace(message string, args ...any) {
	l.write(LogLevelTrace, "VERBOSE", message, args...)
}

func (l *Logger) Fatal(message string, args ...any) {
	l.Error(message, args...)
	os.Exit(1)
}

func (l *Logger) Panic(message string, args ...any) {
	l.Error(message, args...)
	panic(message)
}

func (l *Logger) write(level LogLevel, label, message string, args ...any) {
	if !l.Enabled(level) {
		return
	}
	if _, ignored := ignoredLogContexts[l.context]; ignored {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	labelField := fmt.Sprintf("%5s", label)
	if l.core.color {
		labelField = colorizeLabel(labelField, label)
	}

	contextPart := ""
	if l.context != "" {
		contextPart = " [" + l.context + "]"
		if l.core.color {
			contextPart = colorize(contextPart, ansiYellow)
		}
	}

	line := fmt.Sprintf("%s %s%s %s%s", timestamp, labelField, contextPart, message, formatLogArgs(args...))

	l.core.mu.Lock()
	_, _ = fmt.Fprintln(l.core.writer, line)
	l.core.mu.Unlock()
}

func formatLogArgs(args ...any) string {
	if len(args) == 0 {
		return ""
	}

	fields := make([]string, 0, (len(args)+1)/2)
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			fields = append(fields, fmt.Sprintf("%q: %s", "extra", jsonLogValue(args[i])))
			break
		}

		key, ok := args[i].(string)
		if !ok || strings.TrimSpace(key) == "" {
			fields = append(fields, fmt.Sprintf("%q: %s", "extra", jsonLogValue(args[i])))
			i--
			continue
		}
		fields = append(fields, fmt.Sprintf("%q: %s", key, jsonLogValue(args[i+1])))
	}
	if len(fields) == 0 {
		return ""
	}
	return " {" + strings.Join(fields, ", ") + "}"
}

func jsonLogValue(value any) string {
	switch v := value.(type) {
	case nil:
		return "null"
	case error:
		return quoteJSONString(v.Error())
	case fmt.Stringer:
		return quoteJSONString(v.String())
	case string:
		return quoteJSONString(v)
	default:
		data, err := json.Marshal(v)
		if err == nil {
			return string(data)
		}
		return quoteJSONString(fmt.Sprint(v))
	}
}

func quoteJSONString(value string) string {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%q", value)
	}
	return string(data)
}

const (
	ansiReset   = "\033[0m"
	ansiRed     = "\033[31m"
	ansiYellow  = "\033[33m"
	ansiGreen   = "\033[32m"
	ansiMagenta = "\033[35m"
	ansiCyan    = "\033[36m"
)

func colorizeLabel(value, label string) string {
	switch label {
	case "ERROR":
		return colorize(value, ansiRed)
	case "WARN":
		return colorize(value, ansiYellow)
	case "DEBUG":
		return colorize(value, ansiMagenta)
	case "VERBOSE":
		return colorize(value, ansiCyan)
	default:
		return colorize(value, ansiGreen)
	}
}

func colorize(value, color string) string {
	return color + value + ansiReset
}
