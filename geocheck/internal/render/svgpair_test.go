package render

import (
	"bytes"
	lg "charm.land/lipgloss/v2"
	"regexp"
	"strings"
	"testing"
)

func lineText(l textLine) string {
	var b strings.Builder
	for _, r := range l.runs {
		b.WriteString(r.text)
	}
	return b.String()
}

// contentOf drops rule characters. A section header's rule is shortened to its
// own column by design, so comparing it byte for byte would report an intended
// change as lost content.
func contentOf(l textLine) string {
	return strings.TrimSpace(strings.ReplaceAll(lineText(l), "─", ""))
}

func demoLines(t *testing.T) []textLine {
	t.Helper()
	return ansiLines(renderANSI(DemoReport("test"), nil))
}

// TestGeoTablesSharuARow is the arrangement this code exists for: the two
// short tables sit beside the tall one instead of under it.
func TestGeoTablesShareARow(t *testing.T) {
	paired := pairGeoTables(demoLines(t))

	found := false
	for _, l := range paired {
		text := lineText(l)
		if strings.Contains(text, leftSection) && strings.Contains(text, rightSections[0]) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("%q and %q did not end up on the same row", leftSection, rightSections[0])
	}
}

// TestPairingKeepsEveryLine is the property that matters: moving things must
// not drop or invent a row.
func TestPairingKeepsEveryLine(t *testing.T) {
	lines := demoLines(t)
	paired := pairGeoTables(lines)

	missing := 0
	for _, l := range lines {
		want := contentOf(l)
		if want == "" {
			continue
		}
		found := false
		for _, p := range paired {
			if strings.Contains(contentOf(p), want) {
				found = true
				break
			}
		}
		if !found {
			missing++
			if missing <= 3 {
				t.Errorf("pairing lost %.60q", want)
			}
		}
	}
	if missing > 0 {
		t.Errorf("%d rows went missing", missing)
	}
}

func TestPairingIsShorter(t *testing.T) {
	lines := demoLines(t)
	if got := len(pairGeoTables(lines)); got >= len(lines) {
		t.Errorf("paired to %d rows from %d; it should be shorter", got, len(lines))
	}
}

// TestPairingLeavesOtherReportsAlone covers the runs where the sections are not
// all present — `-g cdn`, or any report with the geolocation half skipped.
func TestPairingLeavesOtherReportsAlone(t *testing.T) {
	r := DemoReport("test")
	r.Geo = nil
	lines := ansiLines(renderANSI(r, nil))

	paired := pairGeoTables(lines)
	if len(paired) != len(lines) {
		t.Fatalf("a report without those sections was rearranged: %d rows became %d",
			len(lines), len(paired))
	}
	for i := range lines {
		if lineText(lines[i]) != lineText(paired[i]) {
			t.Fatalf("row %d changed", i)
		}
	}
}

// TestRulesDoNotWidenThePicture guards the mistake that cost the most: a
// section rule is padded to the terminal's width, so leaving it untrimmed in
// the right-hand column pushed the image half again as wide as it needed to be.
func TestRulesDoNotWidenThePicture(t *testing.T) {
	lines := demoLines(t)

	before := 0
	for _, l := range lines {
		if w := l.width(); w > before {
			before = w
		}
	}
	after := 0
	for _, l := range pairGeoTables(lines) {
		if w := l.width(); w > after {
			after = w
		}
	}
	// Two columns are wider than one, but not by the width of a whole rule.
	if after > before+columnGap+70 {
		t.Errorf("paired layout is %d columns against %d; the rules are not being trimmed",
			after, before)
	}
}

// --- the same arrangement in the terminal -----------------------------------

func printAt(t *testing.T, width int) string {
	t.Helper()
	var buf bytes.Buffer
	o := &Output{W: &buf, Theme: newTheme(lg.LightDark(true)), Color: true, Width: width}
	o.Print(DemoReport("test"))
	return buf.String()
}

func plainLines(s string) []string {
	return strings.Split(regexp.MustCompile(`\x1b\[[0-9;]*m`).ReplaceAllString(s, ""), "\n")
}

func hasPairedRow(s string) bool {
	for _, l := range plainLines(s) {
		if strings.Contains(l, leftSection) && strings.Contains(l, rightSections[0]) {
			return true
		}
	}
	return false
}

// TestTerminalPairsOnlyWhenItFits is the whole contract of doing this on a
// terminal: a narrow one must be left exactly as it was.
func TestTerminalPairsOnlyWhenItFits(t *testing.T) {
	for _, c := range []struct {
		width int
		want  bool
	}{
		{80, false}, {100, false}, {115, false}, {121, false},
		{122, true}, {160, true},
	} {
		if got := hasPairedRow(printAt(t, c.width)); got != c.want {
			t.Errorf("width %d: paired = %v, want %v", c.width, got, c.want)
		}
	}
}

// TestPairingStaysWithinTheTerminal is the failure this must not have: a
// second column pushed past the edge wraps, and a wrapped table is unreadable.
//
// It asserts only about widths where pairing happens. Below that the report is
// already wider than a narrow terminal — the tables are drawn to fit their
// contents and do not shrink — which is long-standing behaviour and not
// something this layout introduced.
func TestPairingStaysWithinTheTerminal(t *testing.T) {
	for _, width := range []int{122, 130, 160, 200} {
		out := printAt(t, width)
		if !hasPairedRow(out) {
			t.Fatalf("width %d should have paired", width)
		}
		for i, l := range plainLines(out) {
			if n := len([]rune(l)); n > width {
				t.Errorf("width %d: line %d is %d columns", width, i, n)
				break
			}
		}
	}
}

// TestNarrowTerminalIsUntouched pins that the stacked path is not merely
// equivalent but identical — it is the same bytes it produced before any of
// this existed.
func TestNarrowTerminalIsUntouched(t *testing.T) {
	var direct bytes.Buffer
	o := &Output{W: &direct, Theme: newTheme(lg.LightDark(true)), Color: true, Width: 100}
	o.printStacked(DemoReport("test"))

	if got := printAt(t, 100); got != direct.String() {
		t.Error("a narrow terminal took the paired path")
	}
}

// TestPairedTerminalKeepsStyling guards the round trip: the paired output is
// rebuilt from parsed runs, so colour and bold have to survive being taken
// apart and put back together.
func TestPairedTerminalKeepsStyling(t *testing.T) {
	out := printAt(t, 130)
	if n := strings.Count(out, "\x1b["); n < 500 {
		t.Errorf("only %d escape sequences survived; styling was lost", n)
	}
	if !strings.Contains(out, "\x1b[1;38;2;") {
		t.Error("bold styling did not survive the round trip")
	}
}
