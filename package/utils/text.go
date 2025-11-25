package utils

import "strings"

func SplitSentences(text string) []string {
	seps := []string{". ", "! ", "? "}
	var result []string
	buf := text

	for {
		minIdx := -1
		sep := ""
		for _, s := range seps {
			if idx := strings.Index(buf, s); idx != -1 && (minIdx == -1 || idx < minIdx) {
				minIdx = idx
				sep = s
			}
		}
		if minIdx == -1 {
			break
		}
		sentence := buf[:minIdx+len(sep)]
		result = append(result, strings.TrimSpace(sentence))
		buf = buf[minIdx+len(sep):]
	}
	if strings.TrimSpace(buf) != "" {
		result = append(result, strings.TrimSpace(buf))
	}
	return result
}
