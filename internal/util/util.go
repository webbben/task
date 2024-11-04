package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Truncate(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength-1] + "-"
	}
	return s
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
