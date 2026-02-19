package helper

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func Center(item any, width int) string {
	return CenterCustom(item, width, "")
}

func Left(item any, width int) string {
	return LeftCustom(item, width, "")
}

func Right(item any, width int) string {
	return RightCustom(item, width, "")
}

func CenterCustom(item any, width int, runes string) string {
	if runes == "" {
		runes = " "
	}

	text := fmt.Sprintf("%v", item)
	if utf8.RuneCountInString(text) >= width {
		return text
	}

	padding := width - utf8.RuneCountInString(text)
	left := padding / 2
	right := padding - left

	return strings.Repeat(runes, left) + text + strings.Repeat(runes, right)
}

func LeftCustom(item any, width int, runes string) string {
	if runes == "" {
		runes = " "
	}

	text := fmt.Sprintf("%v", item)
	if utf8.RuneCountInString(text) >= width {
		return text
	}

	padding := width - utf8.RuneCountInString(text)

	return strings.Repeat(runes, padding) + text
}

func RightCustom(item any, width int, runes string) string {
	if runes == "" {
		runes = " "
	}

	text := fmt.Sprintf("%v", item)
	if utf8.RuneCountInString(text) >= width {
		return text
	}

	padding := width - utf8.RuneCountInString(text)

	return text + strings.Repeat(runes, padding)
}

func Fill(item any, width int) string {
	text := fmt.Sprintf("%v", item)
	if utf8.RuneCountInString(text) >= width {
		return text
	}

	return strings.Repeat(text, width)
}

func RepeatLeft(item any, width int) string {
	return RepeatLeftCustom(item, width, "")
}

func RepeatLeftCustom(item any, width int, runes string) string {
	if runes == "" {
		runes = " "
	}

	text := fmt.Sprintf("%v", item)
	return strings.Repeat(runes, width) + text
}

func RepeatRight(item any, width int) string {
	return RepeatRightCustom(item, width, "")
}

func RepeatRightCustom(item any, width int, runes string) string {
	if runes == "" {
		runes = " "
	}

	text := fmt.Sprintf("%v", item)
	return text + strings.Repeat(runes, width)
}

func NumberToAlpha(n int) string {
	if n <= 0 {
		return "?"
	}

	result := ""

	for n > 0 {
		n--
		remainder := n % 26
		result = string(rune('a'+remainder)) + result
		n = n / 26
	}

	return result
}
