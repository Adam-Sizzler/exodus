package render

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	lg "charm.land/lipgloss/v2"

	"exodus-node/geocheck/internal/detect"
)

//go:embed fonts/jetbrains-mono-regular.woff2 fonts/jetbrains-mono-bold.woff2
var fontFS embed.FS

// The SVG is laid out by column and row rather than by measuring text, which
// only works because the embedded font is monospaced and its advance width is
// known: JetBrains Mono advances 600 units per 1000 em.
const (
	svgFontSize  = 15.0
	svgAdvance   = svgFontSize * 0.6
	svgLineStep  = svgFontSize * 1.34
	svgPadX      = 22.0
	svgPadY      = 20.0
	svgBaseline  = svgFontSize * 1.02 // first baseline below the top padding
	svgRadius    = 10.0
	svgBackplate = "#161b22"
)

// SVG writes the report as a self-contained SVG: fonts embedded, no external
// references, nothing to fetch when it is displayed.
//
// It renders through the ordinary terminal path and then reads the ANSI it
// produced, rather than reimplementing the layout. A second implementation
// would drift from the first the moment either changed, and the point of the
// picture is to show what the terminal showed.
func SVG(w io.Writer, r Report, findings []detect.Finding) error {
	lines := pairGeoTables(ansiLines(renderANSI(r, findings)))

	cols := 0
	for _, line := range lines {
		if n := line.width(); n > cols {
			cols = n
		}
	}
	width := 2*svgPadX + float64(cols)*svgAdvance
	height := 2*svgPadY + float64(len(lines))*svgLineStep

	// The fallbacks are not decoration. Renderers built on librsvg ignore an
	// @font-face with a data: URI, and would otherwise draw the whole report
	// as empty boxes. Since every run is placed by column rather than by
	// measuring text, a substituted font changes glyph shapes but not the
	// layout, so the picture stays readable.
	const family = `JetBrains Mono, ui-monospace, DejaVu Sans Mono, Menlo, Consolas, monospace`

	var b bytes.Buffer
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f" font-family="%s" font-size="%g">`,
		width, height, width, height, family, svgFontSize)

	regular, err := fontFace("fonts/jetbrains-mono-regular.woff2", 400)
	if err != nil {
		return err
	}
	bold, err := fontFace("fonts/jetbrains-mono-bold.woff2", 700)
	if err != nil {
		return err
	}
	// Ligatures are switched off deliberately: JetBrains Mono turns "->" and
	// "--" into single glyphs, which would silently rewrite hostnames and flag
	// names in a picture people read as a transcript.
	fmt.Fprintf(&b, `<style>%s%s text{font-variant-ligatures:none;white-space:pre}</style>`,
		regular, bold)

	fmt.Fprintf(&b, `<rect width="%.0f" height="%.0f" rx="%g" fill="%s"/>`,
		width, height, svgRadius, svgBackplate)

	for row, line := range lines {
		y := svgPadY + svgBaseline + float64(row)*svgLineStep
		for _, run := range line.runs {
			if strings.TrimSpace(run.text) == "" {
				continue
			}
			x := svgPadX + float64(run.col)*svgAdvance
			fmt.Fprintf(&b, `<text x="%.2f" y="%.2f" fill="%s"`, x, y, run.colour())
			if run.bold {
				b.WriteString(` font-weight="700"`)
			}
			b.WriteString(">")
			xmlEscape(&b, run.text)
			b.WriteString(`</text>`)
		}
	}

	b.WriteString(`</svg>`)
	_, err = w.Write(b.Bytes())
	return err
}

// SVGDataURIPrefix is what turns the base64 below into something an <img> or a
// browser address bar accepts directly.
const SVGDataURIPrefix = "data:image/svg+xml;base64,"

// SVGBase64 writes the same document base64-encoded, which is what a chat
// window or a data: URI wants.
func SVGBase64(w io.Writer, r Report, findings []detect.Finding) error {
	var buf bytes.Buffer
	if err := SVG(&buf, r, findings); err != nil {
		return err
	}
	enc := base64.NewEncoder(base64.StdEncoding, w)
	if _, err := enc.Write(buf.Bytes()); err != nil {
		return err
	}
	return enc.Close()
}

func fontFace(path string, weight int) (string, error) {
	data, err := fontFS.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("embedded font %s: %w", path, err)
	}
	return fmt.Sprintf(
		`@font-face{font-family:"JetBrains Mono";font-style:normal;font-weight:%d;src:url(data:font/woff2;base64,%s)format("woff2")}`,
		weight, base64.StdEncoding.EncodeToString(data)), nil
}

// renderANSI produces the report exactly as a colour terminal would receive it.
func svgPalette() palette {
	return palette{
		fg:     lg.Color("#c9d1d9"),
		muted:  lg.Color("#8b949e"),
		accent: lg.Color("#54aeff"),
		good:   lg.Color("#aceebb"),
		warn:   lg.Color("#fae17d"),
		bad:    lg.Color("#ff7b72"),
		info:   lg.Color("#ffb77c"),
		rule:   lg.Color("#30363d"),
	}
}

func renderANSI(r Report, findings []detect.Finding) string {
	var buf bytes.Buffer
	o := &Output{
		W:     &buf,
		Theme: themeFrom(svgPalette()),
		Color: true,
		Width: 115,
	}
	o.PrintFindings(findings)
	o.Print(r)
	return buf.String()
}

func (r textRun) colour() string {
	if !r.hasFG {
		return "#c9d1d9"
	}
	return fmt.Sprintf("#%02x%02x%02x", r.fg[0], r.fg[1], r.fg[2])
}

func xmlEscape(b *bytes.Buffer, s string) {
	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		default:
			b.WriteRune(r)
		}
	}
}
