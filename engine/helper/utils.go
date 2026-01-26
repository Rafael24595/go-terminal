package helper

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func Center(item any, width int) string {
	text := fmt.Sprintf("%v", item)
	if utf8.RuneCountInString(text) >= width {
		return text
	}

	padding := width - utf8.RuneCountInString(text)
	left := padding / 2
	right := padding - left

	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}

func Left(item any, width int) string {
	text := fmt.Sprintf("%v", item)
	if utf8.RuneCountInString(text) >= width {
		return text
	}

	padding := width - utf8.RuneCountInString(text)

	return strings.Repeat(" ", padding) + text
}

func Right(item any, width int) string {
	text := fmt.Sprintf("%v", item)
	if utf8.RuneCountInString(text) >= width {
		return text
	}

	padding := width - utf8.RuneCountInString(text)

	return text + strings.Repeat(" ", padding)
}

func Fill(item any, width int) string {
	text := fmt.Sprintf("%v", item)
	if utf8.RuneCountInString(text) >= width {
		return text
	}

	return strings.Repeat(text, width)
}

func RepeatLeft(item any, width int) string {
	text := fmt.Sprintf("%v", item)
	return strings.Repeat(" ", width) + text
}

func RepeatRight(item any, width int) string {
	text := fmt.Sprintf("%v", item)
	return text + strings.Repeat(" ", width)
}
