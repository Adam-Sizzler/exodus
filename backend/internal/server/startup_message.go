package server

import (
	"strings"
)

const startupTableWidth = 60

func GetStartMessage(version string) string {
	version = strings.TrimSpace(version)
	if version == "" {
		version = "v0.0.0-dev"
	}

	return renderSingleColumnTable(
		"Exodus Subscription Page "+version,
		[]string{"Docs → https://docs.exodus.dev\nCommunity → https://t.me/exodus"},
	)
}

func renderSingleColumnTable(title string, rows []string) string {
	inner := startupTableWidth
	border := strings.Repeat("─", inner+2)

	var b strings.Builder
	b.WriteString("╭")
	b.WriteString(border)
	b.WriteString("╮\n")
	b.WriteString(renderCenteredLine(title, inner))
	b.WriteString("\n├")
	b.WriteString(border)
	b.WriteString("┤\n")

	for _, row := range rows {
		for _, line := range strings.Split(row, "\n") {
			b.WriteString(renderCenteredLine(line, inner))
			b.WriteString("\n")
		}
	}

	b.WriteString("╰")
	b.WriteString(border)
	b.WriteString("╯")
	return b.String()
}

func renderCenteredLine(value string, width int) string {
	value = truncateRunes(strings.TrimSpace(value), width)
	padding := width - runeLen(value)
	left := padding / 2
	right := padding - left
	return "│ " + strings.Repeat(" ", left) + value + strings.Repeat(" ", right) + " │"
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
