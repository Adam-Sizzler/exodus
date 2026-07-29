package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

type actionKey int

const (
	actionKeyNone actionKey = iota
	actionKeyUp
	actionKeyDown
	actionKeyEnter
	actionKeyInterrupt
)

func promptAction() (string, error) {
	actions := rescueActions()
	if len(actions) == 0 {
		return "", errors.New("no rescue actions configured")
	}

	return promptSelect(actions, len(actions)-1)
}

func promptSelect(actions []cliAction, initialIndex int) (string, error) {
	if len(actions) == 0 {
		return "", errors.New("no options configured")
	}

	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		return promptSelectPlain(actions, initialIndex)
	}

	selected := initialIndex
	if selected < 0 || selected >= len(actions) {
		selected = len(actions) - 1
	}

	fmt.Println()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("enable raw terminal mode: %w", err)
	}

	terminalRestored := false
	restoreTerminal := func() {
		if terminalRestored {
			return
		}

		_ = term.Restore(int(os.Stdin.Fd()), oldState)
		terminalRestored = true
	}
	defer restoreTerminal()

	renderActionPrompt(actions, selected)

	for {
		key, err := readActionKey(os.Stdin)
		if err != nil {
			restoreTerminal()
			clearPromptArea(len(actions) + 2)
			return "", fmt.Errorf("read action key: %w", err)
		}

		switch key {
		case actionKeyUp:
			selected--
			if selected < 0 {
				selected = len(actions) - 1
			}
			clearPromptArea(len(actions) + 1)
			renderActionPrompt(actions, selected)
		case actionKeyDown:
			selected++
			if selected >= len(actions) {
				selected = 0
			}
			clearPromptArea(len(actions) + 1)
			renderActionPrompt(actions, selected)
		case actionKeyEnter:
			restoreTerminal()
			clearPromptArea(len(actions) + 2)
			printSelectedAction(actions[selected].Label)
			return actions[selected].Value, nil
		case actionKeyInterrupt:
			restoreTerminal()
			clearPromptArea(len(actions) + 2)
			printStatus("ℹ", "👋 Exiting...")
			os.Exit(0)
		}
	}
}

func promptSelectPlain(actions []cliAction, initialIndex int) (string, error) {
	if initialIndex < 0 || initialIndex >= len(actions) {
		initialIndex = len(actions) - 1
	}

	fmt.Println()
	fmt.Println("Select an action:")
	for index, action := range actions {
		if action.Hint != "" {
			fmt.Printf("%d) %s (%s)\n", index+1, action.Label, action.Hint)
			continue
		}

		fmt.Printf("%d) %s\n", index+1, action.Label)
	}

	fmt.Printf("Enter number [%d]: ", initialIndex+1)

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("read action: %w", err)
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return actions[initialIndex].Value, nil
	}

	selected, err := strconv.Atoi(text)
	if err != nil || selected < 1 || selected > len(actions) {
		return "", fmt.Errorf("invalid action: %s", text)
	}

	selected--
	printSelectedAction(actions[selected].Label)

	return actions[selected].Value, nil
}

func renderActionPrompt(actions []cliAction, selected int) {
	fmt.Printf("%s %s\r\n", ansi("36", "❯"), ansi("36", "Select an action"))

	for index, action := range actions {
		if index == selected {
			fmt.Printf("%s\r\n", formatActiveAction(action))
			continue
		}

		fmt.Printf("%s\r\n", formatInactiveAction(action))
	}
}

func formatActiveAction(action cliAction) string {
	line := fmt.Sprintf("%s %s", ansi("32", "●"), action.Label)
	if action.Hint != "" {
		line += fmt.Sprintf(" %s", ansi("90", fmt.Sprintf("(%s)", action.Hint)))
	}

	return line
}

func formatInactiveAction(action cliAction) string {
	return fmt.Sprintf("%s %s", ansi("90", "○"), ansi("90", action.Label))
}

func readActionKey(reader io.Reader) (actionKey, error) {
	var input [1]byte
	if _, err := reader.Read(input[:]); err != nil {
		return actionKeyNone, err
	}

	switch input[0] {
	case 3:
		return actionKeyInterrupt, nil
	case '\r', '\n':
		return actionKeyEnter, nil
	case 'k', 'K':
		return actionKeyUp, nil
	case 'j', 'J':
		return actionKeyDown, nil
	case 27:
		var sequence [2]byte
		if _, err := io.ReadFull(reader, sequence[:]); err != nil {
			return actionKeyNone, err
		}

		if sequence[0] != '[' {
			return actionKeyNone, nil
		}

		switch sequence[1] {
		case 'A':
			return actionKeyUp, nil
		case 'B':
			return actionKeyDown, nil
		}
	}

	return actionKeyNone, nil
}

func printInfoBox(lines []string) {
	width := 0
	for _, line := range lines {
		if l := len([]rune(line)); l > width {
			width = l
		}
	}
	width += 4

	border := strings.Repeat("─", width)

	fmt.Println()
	fmt.Printf(" ╭%s╮\n", border)
	for _, line := range lines {
		padding := strings.Repeat(" ", width-2-len([]rune(line)))
		fmt.Printf(" │  %s%s │\n", line, padding)
	}
	fmt.Printf(" ╰%s╯\n", border)
	fmt.Println()
}

func promptConfirm(reader *bufio.Reader, label string) (bool, error) {
	for {
		fmt.Printf("? %s [y/N]: ", label)

		text, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("read confirmation: %w", err)
		}

		switch strings.ToLower(strings.TrimSpace(text)) {
		case "y", "yes", "true", "1":
			return true, nil
		case "", "n", "no", "false", "0":
			fmt.Println("❌ Aborted.")
			return false, nil
		default:
			fmt.Println("Please answer yes or no.")
		}
	}
}

func printCLIBox(title string) {
	width := len([]rune(title)) + 4
	line := strings.Repeat("─", width)
	empty := strings.Repeat(" ", width)

	fmt.Println()
	fmt.Printf(" ╭%s╮\n", line)
	fmt.Printf(" │%s│\n", empty)
	fmt.Printf(" │  %s  │\n", title)
	fmt.Printf(" │%s│\n", empty)
	fmt.Printf(" ╰%s╯\n", line)
	fmt.Println()
}

func printSelectedAction(label string) {
	fmt.Println()
	fmt.Printf("%s Select an action\n", ansi("32", "✔"))
	fmt.Println(label)
}

func printStatus(symbol string, message string) {
	now := time.Now().Format("3:04:05 PM")
	coloredNow := ansi("90", now)
	width := terminalWidth()

	plainLength := len([]rune(symbol)) + 1 + len([]rune(message)) + 1 + len([]rune(now))

	if plainLength >= width {
		fmt.Printf("%s %s %s\n", symbol, message, coloredNow)
		return
	}

	padding := strings.Repeat(" ", width-plainLength)
	fmt.Printf("%s %s%s%s\n", symbol, message, padding, coloredNow)
}

func terminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil && width > 0 {
		return width
	}

	columns := strings.TrimSpace(os.Getenv("COLUMNS"))
	if columns != "" {
		parsedWidth, err := strconv.Atoi(columns)
		if err == nil && parsedWidth > 0 {
			return parsedWidth
		}
	}

	return 120
}

func clearPromptArea(lines int) {
	if lines <= 0 {
		return
	}

	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return
	}

	fmt.Printf("\x1b[%dA", lines)
	fmt.Print("\x1b[J")
}

func ansi(code string, value string) string {
	return "\x1b[" + code + "m" + value + "\x1b[0m"
}
