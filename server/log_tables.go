package server

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"exodus-node/config"
	"exodus-node/constant"
)

func GetStartMessage(cfg *config.NodeConfig) string {
	port := config.DefaultNodeGRPCPort
	if cfg != nil {
		port = cfg.Exodus.GrpcPort
	}

	return renderGroupedLogTable(
		"Exodus Node "+constant.Version,
		[][]string{
			{
				"Docs → https://docs.exodus.dev",
				"Community → https://t.me/exodus",
			},
			{
				"API Port → " + strconv.Itoa(port),
				"Internal Ports → " + strconv.Itoa(config.FixedCoreAPIGRPCPort),
			},
			{
				"Sing-box Core → v" + detectManagedCoreVersion(),
				"Sing-box Path → /usr/local/bin/sing-box",
			},
			{fmt.Sprintf("%dC, %s, %s", runtime.NumCPU(), detectCPUModelForLogs(), formatIECBytesForLogs(detectTotalRAMForLogs()))},
			{"Kernel → " + strings.TrimSpace(runCommandForLogs("uname", "-r"))},
			wrapPrefixedLogLine("Interfaces → ", strings.Join(detectNetworkInterfacesForLogs(), ", "), logTableWidth-4),
		},
	)
}

func renderCoreStartedMessage(processState string) string {
	return renderPlainLogTable(
		"Sing-box started successfully",
		[]string{
			"Version → " + detectManagedCoreVersion(),
			"Process State → " + processState,
			"Started At → " + time.Now().Format(time.RFC3339),
			"Config → " + config.FixedSingboxConfigPath,
		},
	)
}

func renderCoreFailedMessage(processState, err string) string {
	return renderPlainLogTable(
		"Sing-box failed to start",
		append([]string{
			"Version → " + detectManagedCoreVersion(),
			"Process State → " + processState,
		}, wrapPrefixedLogLine("Error → ", err, logTableWidth-4)...),
	)
}

const logTableWidth = 80

func renderGroupedLogTable(title string, groups [][]string) string {
	width := logTableWidth
	borderTop := "╭" + strings.Repeat("─", width-2) + "╮"
	borderMid := "├" + strings.Repeat("─", width-2) + "┤"
	borderBottom := "╰" + strings.Repeat("─", width-2) + "╯"

	var b strings.Builder
	b.WriteString(borderTop)
	b.WriteByte('\n')
	b.WriteString(renderCenteredLine(title, width))
	b.WriteByte('\n')
	b.WriteString(borderMid)

	for groupIndex, group := range groups {
		if groupIndex > 0 {
			b.WriteByte('\n')
			b.WriteString(renderDashedLine(width))
		}
		for _, line := range group {
			for _, wrapped := range wrapLogLine(line, width-4) {
				b.WriteByte('\n')
				b.WriteString(renderCenteredLine(wrapped, width))
			}
		}
	}

	b.WriteByte('\n')
	b.WriteString(borderBottom)
	return b.String()
}

func renderPlainLogTable(title string, lines []string) string {
	width := logTableWidth
	borderTop := "╭" + strings.Repeat("─", width-2) + "╮"
	borderMid := "├" + strings.Repeat("─", width-2) + "┤"
	borderBottom := "╰" + strings.Repeat("─", width-2) + "╯"

	var b strings.Builder
	b.WriteString(borderTop)
	b.WriteByte('\n')
	b.WriteString(renderCenteredLine(title, width))
	b.WriteByte('\n')
	b.WriteString(borderMid)
	for _, line := range lines {
		for _, wrapped := range wrapLogLine(line, width-4) {
			b.WriteByte('\n')
			b.WriteString(renderCenteredLine(wrapped, width))
		}
	}
	b.WriteByte('\n')
	b.WriteString(borderBottom)
	return b.String()
}

func renderDashedLine(width int) string {
	return "│" + strings.Repeat("-", width-2) + "│"
}

func renderCenteredLine(value string, width int) string {
	inner := width - 2
	value = strings.TrimSpace(value)
	if logCellWidth(value) > inner {
		value = truncateLogCell(value, inner)
	}
	used := logCellWidth(value)
	left := (inner - used) / 2
	right := inner - used - left
	return "│" + strings.Repeat(" ", left) + value + strings.Repeat(" ", right) + "│"
}

func wrapPrefixedLogLine(prefix, value string, max int) []string {
	prefix = strings.TrimSpace(prefix) + " "
	continuation := strings.Repeat(" ", logCellWidth(prefix))
	wrapped := wrapLogLine(prefix+strings.TrimSpace(value), max)
	for i := 1; i < len(wrapped); i++ {
		wrapped[i] = continuation + strings.TrimSpace(wrapped[i])
	}
	return wrapped
}

func wrapLogLine(value string, max int) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return []string{""}
	}
	if logCellWidth(value) <= max {
		return []string{value}
	}

	parts := strings.Split(value, ", ")
	if len(parts) == 1 {
		return wrapByWords(value, max)
	}

	lines := make([]string, 0, len(parts))
	current := ""
	for _, part := range parts {
		part = strings.TrimSpace(part)
		candidate := part
		if current != "" {
			candidate = current + ", " + part
		}
		if current != "" && logCellWidth(candidate) > max {
			lines = append(lines, current)
			current = part
			continue
		}
		current = candidate
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func wrapByWords(value string, max int) []string {
	words := strings.Fields(value)
	lines := make([]string, 0, len(words))
	current := ""
	for _, word := range words {
		candidate := word
		if current != "" {
			candidate = current + " " + word
		}
		if current != "" && logCellWidth(candidate) > max {
			lines = append(lines, current)
			current = word
			continue
		}
		current = candidate
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func truncateLogCell(value string, max int) string {
	value = strings.TrimSpace(value)
	if logCellWidth(value) <= max {
		return value
	}
	if max <= 1 {
		return string([]rune(value)[:max])
	}

	runes := []rune(value)
	if len(runes) > max-1 {
		runes = runes[:max-1]
	}
	return string(runes) + "…"
}

func logCellWidth(value string) int {
	return len([]rune(value))
}

func detectManagedCoreVersion() string {
	out := runCommandForLogs("/usr/local/bin/sing-box", "version")
	if out == "" {
		return "N/A"
	}
	line := strings.TrimSpace(strings.SplitN(out, "\n", 2)[0])
	if idx := strings.LastIndex(line, " "); idx > 0 && idx+1 < len(line) {
		return strings.TrimSpace(line[idx+1:])
	}
	return line
}

func runCommandForLogs(name string, args ...string) string {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func detectCPUModelForLogs() string {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return runtime.GOARCH
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(strings.ToLower(line), "model name") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return runtime.GOARCH
}

func detectTotalRAMForLogs() uint64 {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				kb, err := strconv.ParseUint(fields[1], 10, 64)
				if err == nil {
					return kb * 1024
				}
			}
		}
	}
	return 0
}

func formatIECBytesForLogs(value uint64) string {
	if value == 0 {
		return "N/A"
	}
	units := []string{"B", "KiB", "MiB", "GiB", "TiB"}
	floatValue := float64(value)
	idx := 0
	for floatValue >= 1024 && idx < len(units)-1 {
		floatValue /= 1024
		idx++
	}
	return fmt.Sprintf("%.1f %s", floatValue, units[idx])
}

func detectNetworkInterfacesForLogs() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return []string{"unknown"}
	}
	result := make([]string, 0, len(interfaces))
	for _, iface := range interfaces {
		if iface.Name == "lo" {
			continue
		}
		result = append(result, iface.Name)
	}
	if len(result) == 0 {
		return []string{"unknown"}
	}
	return result
}

func detectDefaultNetworkInterfaceForLogs() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "unknown"
	}
	for _, iface := range interfaces {
		if iface.Name == "lo" || iface.Flags&net.FlagUp == 0 {
			continue
		}
		return iface.Name
	}
	return "unknown"
}
