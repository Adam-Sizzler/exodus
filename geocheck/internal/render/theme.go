// Package render turns check results into terminal output and JSON.
package render

import (
	"image/color"
	"io"
	"os"
	"strconv"
	"strings"

	lg "charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
	xterm "github.com/charmbracelet/x/term"
)

// Theme holds every style the report uses.
type Theme struct {
	Title     lg.Style
	Subtitle  lg.Style
	Section   lg.Style
	Label     lg.Style
	Value     lg.Style
	Muted     lg.Style
	Good      lg.Style
	Warn      lg.Style
	Bad       lg.Style
	Info      lg.Style
	Accent    lg.Style
	TableHead lg.Style
	TableCell lg.Style
	Border    lg.Style
}

// Output bundles a writer with the theme appropriate for it.
type Output struct {
	W     io.Writer
	Theme Theme
	// Color reports whether the destination can show colour at all.
	Color bool
	// Width is the usable terminal width, or a sensible default when piped.
	Width int
}

// NewOutput inspects the destination and builds a matching writer and theme.
// Piped output, NO_COLOR and dumb terminals all degrade to plain text.
func NewOutput(f *os.File) *Output {
	env := os.Environ()
	profile := colorprofile.Detect(f, env)

	// colorprofile lets CLICOLOR_FORCE override NO_COLOR; NO_COLOR must win.
	if os.Getenv("NO_COLOR") != "" {
		profile = colorprofile.Ascii
	}

	isTTY := xterm.IsTerminal(f.Fd())
	colored := profile > colorprofile.Ascii

	isDark := backgroundIsDark()

	width := 100
	if isTTY {
		if w, _, err := xterm.GetSize(f.Fd()); err == nil && w > 40 {
			width = w
		}
	}

	w := colorprofile.NewWriter(f, env)
	w.Profile = profile

	return &Output{W: w, Theme: newTheme(lg.LightDark(isDark)), Color: colored, Width: width}
}

// backgroundIsDark decides which palette to use without ever writing to the
// terminal.
//
// The tempting alternative is to ask: emit an OSC 11 query and read the reply.
// That is unsafe here. Under `docker run -t` without `-i` — which is what a
// container gets when stdout is a terminal but stdin is a pipe — stdin *looks*
// like a terminal while carrying nothing, so the query goes out, the terminal
// answers, and nobody consumes the answer. It lands on the user's screen as
// `^[]11;rgb:1a1a/1b1b/2626^[\`. The same happens over ssh sessions and inside
// multiplexers that do not proxy the reply.
//
// So: read COLORFGBG when the terminal volunteers it, and otherwise assume
// dark, which is what the overwhelming majority of terminals, CI logs and
// containers are. Getting this wrong costs a little contrast; emitting an
// unanswerable escape costs the user a corrupted prompt.
func backgroundIsDark() bool {
	// COLORFGBG is "foreground;background" with ANSI colour indices, e.g.
	// "15;0" for white on black. Some terminals insert a third field.
	fields := strings.Split(os.Getenv("COLORFGBG"), ";")
	if len(fields) < 2 {
		return true
	}
	bg, err := strconv.Atoi(strings.TrimSpace(fields[len(fields)-1]))
	if err != nil {
		return true
	}
	// 0-6 and 8 are the dark half of the ANSI palette; 7 and 9-15 are light.
	return bg <= 6 || bg == 8
}

type palette struct {
	fg, muted, accent, good, warn, bad, info, rule color.Color
}

func newTheme(ld lg.LightDarkFunc) Theme {
	return themeFrom(palette{
		fg:     ld(lg.Color("#1a1a1a"), lg.Color("#e6e6e6")),
		muted:  ld(lg.Color("#6b6b6b"), lg.Color("#8a8a8a")),
		accent: ld(lg.Color("#7b3fe4"), lg.Color("#b19cff")),
		good:   ld(lg.Color("#1f8a44"), lg.Color("#5ee08a")),
		warn:   ld(lg.Color("#a86800"), lg.Color("#f5c453")),
		bad:    ld(lg.Color("#c02020"), lg.Color("#ff7b72")),
		info:   ld(lg.Color("#0a6b96"), lg.Color("#6fd3f7")),
		rule:   ld(lg.Color("#c8c8c8"), lg.Color("#454545")),
	})
}

func themeFrom(p palette) Theme {
	fg, muted, accent := p.fg, p.muted, p.accent
	good, warn, bad, info, rule := p.good, p.warn, p.bad, p.info, p.rule

	return Theme{
		Title:     lg.NewStyle().Bold(true).Foreground(accent),
		Subtitle:  lg.NewStyle().Foreground(muted),
		Section:   lg.NewStyle().Bold(true).Foreground(accent),
		Label:     lg.NewStyle().Foreground(muted),
		Value:     lg.NewStyle().Bold(true).Foreground(fg),
		Muted:     lg.NewStyle().Foreground(muted),
		Good:      lg.NewStyle().Foreground(good),
		Warn:      lg.NewStyle().Foreground(warn),
		Bad:       lg.NewStyle().Foreground(bad),
		Info:      lg.NewStyle().Foreground(info),
		Accent:    lg.NewStyle().Foreground(accent),
		TableHead: lg.NewStyle().Bold(true).Foreground(accent).Padding(0, 1),
		TableCell: lg.NewStyle().Foreground(fg).Padding(0, 1),
		Border:    lg.NewStyle().Foreground(rule),
	}
}
