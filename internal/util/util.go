package util

import (
	"bufio"
	"fmt"
	"os"
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
