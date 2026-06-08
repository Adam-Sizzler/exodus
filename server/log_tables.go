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

	return renderLogTable(
		"Exodus Node "+constant.Version,
		[][]string{
			{"Docs", "https://docs.exodus.dev"},
			{"Community", "https://t.me/exodus"},
			{"API Port", strconv.Itoa(port)},
			{"Internal Ports", strconv.Itoa(config.FixedCoreAPIGRPCPort)},
			{"Sing-box Core", "v" + detectManagedCoreVersion()},
			{"Sing-box Path", "/usr/local/bin/sing-box"},
			{"System", fmt.Sprintf("%dC, %s, %s", runtime.NumCPU(), detectCPUModelForLogs(), formatIECBytesForLogs(detectTotalRAMForLogs()))},
			{"Kernel", strings.TrimSpace(runCommandForLogs("uname", "-r"))},
			{"Network Interfaces", strings.Join(detectNetworkInterfacesForLogs(), ", ")},
		},
	)
}

func renderCoreStartedMessage(processState string) string {
	return renderLogTable(
		"Sing-box started successfully",
		[][]string{
			{"Version", detectManagedCoreVersion()},
			{"Process State", processState},
			{"Started At", time.Now().Format(time.RFC3339)},
			{"Config", config.FixedSingboxConfigPath},
		},
	)
}

func renderCoreFailedMessage(processState, err string) string {
	return renderLogTable(
		"Sing-box failed to start",
		[][]string{
			{"Version", detectManagedCoreVersion()},
			{"Process State", processState},
			{"Error", err},
		},
	)
}

const (
	logTableWidth           = 80
	logTableLeftColumnWidth = 20
)

func renderLogTable(title string, rows [][]string) string {
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
	for _, row := range rows {
		if len(row) < 2 {
			continue
		}
		left := row[0]
		lines := strings.Split(row[1], "\n")
		for i, line := range lines {
			label := ""
			if i == 0 {
				label = left
			}
			b.WriteByte('\n')
			b.WriteString(renderTwoColumnLine(label, line, width))
		}
	}
	b.WriteByte('\n')
	b.WriteString(borderBottom)
	return b.String()
}

func renderCenteredLine(value string, width int) string {
	inner := width - 2
	value = truncateLogCell(value, inner)
	used := logCellWidth(value)
	left := (inner - used) / 2
	right := inner - used - left
	return "│" + strings.Repeat(" ", left) + value + strings.Repeat(" ", right) + "│"
}

func renderTwoColumnLine(left, right string, width int) string {
	inner := width - 2
	leftWidth := logTableLeftColumnWidth
	rightWidth := inner - leftWidth - 5
	return "│ " + renderLogCell(left, leftWidth) + " │ " + renderLogCell(right, rightWidth) + " │"
}

func renderLogCell(value string, width int) string {
	value = truncateLogCell(value, width)
	padding := width - logCellWidth(value)
	if padding < 0 {
		padding = 0
	}
	return value + strings.Repeat(" ", padding)
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
