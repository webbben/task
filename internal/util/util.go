package util

func Truncate(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength-1] + "-"
	}
	return s
}
