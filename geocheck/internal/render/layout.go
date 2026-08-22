package render

// Layout is the step between drawing the report and writing it out. Both the
// terminal and the SVG go through it, which is why it lives apart from either:
// the pairing below is not an SVG feature the terminal borrows, it is the
// layout, and the two writers differ only in how they draw the result.
//
// The report is rendered once as ANSI text, parsed back into positioned runs,
// rearranged here, and then re-emitted as ANSI or drawn as SVG.

import (
	"strconv"
	"strings"
)

// textRun is a stretch of text sharing one style, starting at a column.
type textRun struct {
	col  int
	text string
	bold bool
	fg   [3]uint8
	// hasFG distinguishes "no colour set" from "black", which matters because
	// the default has to come from the theme rather than from zero values.
	hasFG bool
}

type textLine struct {
	runs []textRun
}

func (l textLine) width() int {
	n := 0
	for _, r := range l.runs {
		if end := r.col + len([]rune(r.text)); end > n {
			n = end
		}
	}
	return n
}

// ansiLines turns a rendered report into positioned, styled runs.
//
// Only three sequences ever reach here — reset, a truecolor foreground, and
// the same with bold — because that is all the theme emits. Anything else is
// skipped rather than drawn, so an unrecognised sequence loses styling instead
// of appearing as literal text in the output.
func ansiLines(s string) []textLine {
	var (
		lines []textLine
		cur   textLine
		buf   strings.Builder
		col   int
		state textRun
	)

	flush := func() {
		if buf.Len() == 0 {
			return
		}
		run := state
		run.col = col
		run.text = buf.String()
		cur.runs = append(cur.runs, run)
		col += len([]rune(run.text))
		buf.Reset()
	}

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		switch {
		case runes[i] == '\n':
			flush()
			lines = append(lines, cur)
			cur, col = textLine{}, 0

		case runes[i] == 0x1b && i+1 < len(runes) && runes[i+1] == '[':
			// Scan to the final byte of the CSI sequence. Parameter and
			// intermediate bytes are 0x20-0x3F, the final byte 0x40-0x7E.
			// The whole sequence is consumed even when it is not one we act
			// on: skipping only the escape would spill the parameters into
			// the picture as text.
			end := i + 2
			for end < len(runes) && runes[end] >= 0x20 && runes[end] <= 0x3F {
				end++
			}
			flush()
			if end < len(runes) && runes[end] == 'm' {
				state = applySGR(state, string(runes[i+2:end]))
			}
			i = end

		default:
			buf.WriteRune(runes[i])
		}
	}
	flush()
	if len(cur.runs) > 0 {
		lines = append(lines, cur)
	}

	// A trailing blank line adds height and nothing else.
	for len(lines) > 0 && len(lines[len(lines)-1].runs) == 0 {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func applySGR(state textRun, params string) textRun {
	if params == "" || params == "0" {
		return textRun{}
	}
	fields := strings.Split(params, ";")
	for i := 0; i < len(fields); i++ {
		switch fields[i] {
		case "0":
			state = textRun{}
		case "1":
			state.bold = true
		case "22":
			state.bold = false
		case "39":
			state.hasFG = false
		case "38":
			// 38;2;r;g;b is the only form the theme produces.
			if i+4 < len(fields) && fields[i+1] == "2" {
				r, _ := strconv.Atoi(fields[i+2])
				g, _ := strconv.Atoi(fields[i+3])
				b, _ := strconv.Atoi(fields[i+4])
				state.fg = [3]uint8{clamp8(r), clamp8(g), clamp8(b)}
				state.hasFG = true
				i += 4
			}
		}
	}
	return state
}

func clamp8(v int) uint8 {
	switch {
	case v < 0:
		return 0
	case v > 255:
		return 255
	}
	return uint8(v)
}

// The three geolocation tables are narrow and stacking them wastes most of the
// width. This is the one rearrangement the report makes, in the terminal and
// in the picture alike, and it is spelled out rather than derived: the
// services table is the tall one, so the two shorter tables sit beside it in a
// column of their own.
//
//	Popular services │ GeoIP databases
//	                 │ CDN edge location
const (
	leftSection = "Popular services"
	// columnGap is the space between the two columns, in characters.
	columnGap = 3
)

// rightSections stack in the second column, in this order.
var rightSections = []string{"GeoIP databases", "CDN edge location"}

// pairGeoTables moves the two named sections into a column beside the first.
// Anything it cannot find is left exactly where it was, so a report missing
// those sections — `-g cdn`, say — renders unchanged.
// pairGeoTablesWithin is pairGeoTables with a width budget. It reports false
// when the two columns would not fit, which is how the terminal decides
// between this layout and the stacked one.
func pairGeoTablesWithin(lines []textLine, width int) ([]textLine, bool) {
	out := pairGeoTables(lines)
	if len(out) == len(lines) {
		return nil, false // the sections were not there, or not adjacent
	}
	for _, l := range out {
		if l.width() > width {
			return nil, false
		}
	}
	return out, true
}

func pairGeoTables(lines []textLine) []textLine {
	bounds := sectionBounds(lines)

	left, ok := bounds[leftSection]
	if !ok {
		return lines
	}
	var right [][2]int
	for _, title := range rightSections {
		b, ok := bounds[title]
		if !ok {
			return lines
		}
		right = append(right, b)
	}

	// Only the contiguous run left,right... can be rearranged without changing
	// the order anything else is read in.
	prev := left
	for _, b := range right {
		if b[0] != prev[1] {
			return lines
		}
		prev = b
	}
	end := prev[1]

	leftLines := lines[left[0]:left[1]]
	var rightLines []textLine
	for _, b := range right {
		rightLines = append(rightLines, lines[b[0]:b[1]]...)
	}

	out := make([]textLine, 0, len(lines))
	out = append(out, lines[:left[0]]...)
	out = append(out, twoColumns(leftLines, rightLines)...)
	return append(out, lines[end:]...)
}

// twoColumns lays one set of lines to the right of another.
func twoColumns(left, right []textLine) []textLine {
	leftWidth := contentWidth(left)
	rightWidth := contentWidth(right)
	offset := leftWidth + columnGap

	height := len(left)
	if len(right) > height {
		height = len(right)
	}

	out := make([]textLine, height)
	for i, l := range left {
		out[i].runs = append(out[i].runs, trimRule(l, leftWidth)...)
	}
	for i, l := range right {
		for _, run := range trimRule(l, rightWidth) {
			run.col += offset
			out[i].runs = append(out[i].runs, run)
		}
	}
	return out
}

// contentWidth measures a column ignoring section rules, which are padded to
// the full terminal width and would otherwise decide the answer.
func contentWidth(lines []textLine) int {
	width := 0
	for _, l := range lines {
		if isSectionHeader(l) {
			continue
		}
		if w := l.width(); w > width {
			width = w
		}
	}
	return width
}

// trimRule shortens a section header's rule so it stops at the width of its own
// column instead of running to the width of the whole terminal report.
func trimRule(l textLine, width int) []textRun {
	runs := make([]textRun, 0, len(l.runs))
	for _, r := range l.runs {
		if width > 0 && isSectionHeader(l) && isRule(r.text) {
			n := width - r.col
			if n <= 0 {
				continue
			}
			r.text = strings.Repeat("─", n)
		}
		runs = append(runs, r)
	}
	return runs
}

func isRule(s string) bool {
	t := strings.TrimSpace(s)
	return len(t) >= 3 && strings.Trim(t, "─") == ""
}

// isSectionHeader recognises the line `section` writes: a bold title beside a
// rule. Table borders use the same rule character but carry no bold text,
// which is what tells the two apart.
func isSectionHeader(l textLine) bool {
	bold, rule := false, false
	for _, r := range l.runs {
		if r.bold && strings.TrimSpace(r.text) != "" {
			bold = true
		}
		if isRule(r.text) {
			rule = true
		}
	}
	return bold && rule
}

// sectionBounds maps a section title to the half-open line range it occupies.
func sectionBounds(lines []textLine) map[string][2]int {
	out := map[string][2]int{}
	title, start := "", -1
	for i, l := range lines {
		if !isSectionHeader(l) {
			continue
		}
		if start >= 0 {
			out[title] = [2]int{start, i}
		}
		title, start = headerTitle(l), i
	}
	if start >= 0 {
		out[title] = [2]int{start, len(lines)}
	}
	return out
}

// headerTitle is the bold text a section header opens with.
func headerTitle(l textLine) string {
	for _, r := range l.runs {
		if r.bold {
			if t := strings.TrimSpace(r.text); t != "" {
				return t
			}
		}
	}
	return ""
}

// ansiText turns parsed lines back into the escape sequences a terminal reads.
// Only a truecolor foreground and bold are reproduced, which is all the parser
// records and all the theme emits.
func ansiText(lines []textLine) string {
	var b strings.Builder
	for _, l := range lines {
		col := 0
		for _, run := range l.runs {
			if run.col > col {
				b.WriteString(strings.Repeat(" ", run.col-col))
				col = run.col
			}
			styled := run.bold || run.hasFG
			if styled {
				b.WriteString("\x1b[")
				if run.bold {
					b.WriteString("1;")
				}
				b.WriteString("38;2;")
				b.WriteString(strconv.Itoa(int(run.fg[0])))
				b.WriteString(";")
				b.WriteString(strconv.Itoa(int(run.fg[1])))
				b.WriteString(";")
				b.WriteString(strconv.Itoa(int(run.fg[2])))
				b.WriteString("m")
			}
			b.WriteString(run.text)
			if styled {
				b.WriteString("\x1b[m")
			}
			col += len([]rune(run.text))
		}
		b.WriteString("\n")
	}
	return b.String()
}
