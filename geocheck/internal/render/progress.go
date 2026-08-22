package render

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	xterm "github.com/charmbracelet/x/term"
)

// Progress shows a single self-updating status line while checks run. It draws
// only to a real terminal, so piped and CI output stays clean.
type Progress struct {
	w       *os.File
	enabled bool

	mu      sync.Mutex
	text    string
	frame   int
	lastLen int
	done    chan struct{}
	once    sync.Once
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// NewProgress returns a progress line, active only when enabled and attached
// to a terminal.
func NewProgress(w *os.File, enabled bool) *Progress {
	p := &Progress{
		w:       w,
		enabled: enabled && xterm.IsTerminal(w.Fd()),
		done:    make(chan struct{}),
	}
	if p.enabled {
		go p.animate()
	}
	return p
}

// Set replaces the status text.
func (p *Progress) Set(text string) {
	if p == nil || !p.enabled {
		return
	}
	p.mu.Lock()
	p.text = text
	p.mu.Unlock()
}

// Stop clears the line and halts the animation.
func (p *Progress) Stop() {
	if p == nil || !p.enabled {
		return
	}
	p.once.Do(func() {
		close(p.done)
		p.mu.Lock()
		defer p.mu.Unlock()
		p.clear()
	})
}

func (p *Progress) animate() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-p.done:
			return
		case <-ticker.C:
			p.draw()
		}
	}
}

func (p *Progress) draw() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.text == "" {
		return
	}
	p.frame = (p.frame + 1) % len(spinnerFrames)
	line := fmt.Sprintf("  %s %s", spinnerFrames[p.frame], p.text)
	if pad := p.lastLen - len(line); pad > 0 {
		line += strings.Repeat(" ", pad)
	}
	p.lastLen = len(line)
	fmt.Fprintf(p.w, "\r%s", line)
}

func (p *Progress) clear() {
	if p.lastLen == 0 {
		return
	}
	fmt.Fprintf(p.w, "\r%s\r", strings.Repeat(" ", p.lastLen))
	p.lastLen = 0
}
