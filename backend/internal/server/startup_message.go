package server

import (
	"strconv"
	"strings"

	"github.com/exodus/subscription-page/backend/internal/config"
)

const startupTableWidth = 80

func GetStartMessage(version string, cfg config.Config) string {
	version = strings.TrimSpace(version)
	if version == "" {
		version = "v0.0.0-dev"
	}

	return renderGroupedTable(
		"Exodus Subscription Page "+version,
		[][]string{
			{
				"Docs → https://docs.exodus.dev",
				"Community → https://t.me/exodus",
			},
			{
				"HTTP Port → " + strings.TrimSpace(cfg.AppPort),
				"gRPC Port → " + strconv.Itoa(cfg.GRPCPort),
				"Path Prefix → " + config.DisplayPrefix(cfg.SubPath),
			},
		},
	)
}

func renderGroupedTable(title string, groups [][]string) string {
	width := startupTableWidth
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

func renderDashedLine(width int) string {
	return "│" + strings.Repeat("-", width-2) + "│"
}

func renderCenteredLine(value string, width int) string {
	inner := width - 2
	value = truncateRunes(strings.TrimSpace(value), inner)
	padding := inner - runeLen(value)
	left := padding / 2
	right := padding - left
	return "│" + strings.Repeat(" ", left) + value + strings.Repeat(" ", right) + "│"
}

func truncateRunes(value string, width int) string {
	if width <= 0 || runeLen(value) <= width {
		return value
	}
	runes := []rune(value)
	if width == 1 {
		return "…"
	}
	return string(runes[:width-1]) + "…"
}

func runeLen(value string) int {
	return len([]rune(value))
}

func wrapLogLine(value string, max int) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return []string{""}
	}

	words := strings.Fields(value)
	if len(words) == 0 {
		return []string{value}
	}

	lines := make([]string, 0, 2)
	current := ""
	for _, word := range words {
		if runeLen(word) > max {
			if strings.TrimSpace(current) != "" {
				lines = append(lines, current)
				current = ""
			}
			lines = append(lines, splitLongWord(word, max)...)
			continue
		}
		candidate := word
		if current != "" {
			candidate = current + " " + word
		}
		if runeLen(candidate) > max {
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

func splitLongWord(value string, max int) []string {
	if max <= 0 || runeLen(value) <= max {
		return []string{value}
	}
	runes := []rune(value)
	parts := make([]string, 0, (len(runes)/max)+1)
	for len(runes) > max {
		parts = append(parts, string(runes[:max]))
		runes = runes[max:]
	}
	if len(runes) > 0 {
		parts = append(parts, string(runes))
	}
	return parts
}
