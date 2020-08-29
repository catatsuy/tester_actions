package util

import (
	"strings"
	"unicode/utf8"
)

func AutoDetectJP(input string) bool {
	ratio := float64(len(input)) / float64(utf8.RuneCountInString(input))

	return ratio > 1.5
}

func TrimUnnecessary(input string) string {
	strs := strings.Split(input, "\n")

	newStrs := make([]string, 0, len(strs))
	for _, s := range strs {
		tmp := strings.TrimLeft(s, " *#/\t")
		if tmp == "" {
			tmp = "\n"
		}
		newStrs = append(newStrs, tmp)
	}

	return strings.Join(newStrs, " ")
}
