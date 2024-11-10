package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func Truncate(s string, maxLength int) string {
	ansiCodes := extractANSI(s)

	// ansi codes can make the string appear longer than it actually is visibly
	// so handle ansi encoded strings separately
	if len(ansiCodes) > 0 {
		stripped := StripAnsi(s)
		if len(stripped) > maxLength {
			stripped = stripped[:maxLength-3] + "..."
		}
		return ansiCodes[0] + stripped + ansiCodes[1]
	}

	// for normal strings, just check it as it is
	if len(s) > maxLength {
		return s[:maxLength-3] + "..."
	}
	return s
}

func StripAnsi(str string) string {
	return ansiEscape.ReplaceAllString(str, "")
}

func extractANSI(input string) []string {
	// Find all ANSI sequences in the input string
	return ansiEscape.FindAllString(input, -1)
}

func Confirm(prompt string) bool {
	fmt.Print(prompt + " [Y/N]")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	return input == "y" || input == "yes"
}

func OpenEditor() string {
	// create temp file
	temp, err := os.CreateTemp("", "note-*.txt")
	if err != nil {
		log.Fatal("Failed to create temp file:", err)
	}
	defer os.Remove(temp.Name())

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	// open editor
	cmd := exec.Command(editor, temp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal("Failed to run editor:", err)
	}
	content, err := os.ReadFile(temp.Name())
	if err != nil {
		log.Fatal("failed to read temp file:", err)
	}

	return strings.TrimSpace(string(content))
}

// RoundDateDown returns the earliest time in the same day as the given time
func RoundDateDown(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

// RoundDateUp returns the latest time in the same day as the given time
func RoundDateUp(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())
}
