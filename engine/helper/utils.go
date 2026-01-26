package helper

import (
	"fmt"
	"strings"
)

func Center(item any, width int) string {
	text := fmt.Sprintf("%v", item)
	if len(text) >= width {
		return text
	}

	padding := width - len(text)
	left := padding / 2
	right := padding - left

	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}

func Left(item any, width int) string {
	text := fmt.Sprintf("%v", item)
	if len(text) >= width {
		return text
	}

	padding := width - len(text)

	return strings.Repeat(" ", padding) + text
}

func Right(item any, width int) string {
	text := fmt.Sprintf("%v", item)
	if len(text) >= width {
		return text
	}

	padding := width - len(text)

	return text + strings.Repeat(" ", padding)
}

func Fill(item any, width int) string {
	text := fmt.Sprintf("%v", item)
	if len(text) >= width {
		return text
	}

	return strings.Repeat(text, width)
}
