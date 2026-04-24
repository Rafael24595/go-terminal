package helper

import (
	"fmt"
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
)

type TextLayoutOpts struct {
	LogicalSize int
	Runes       string
}

func fixDirectionOps(text string, opts TextLayoutOpts) TextLayoutOpts {
	if opts.LogicalSize == 0 {
		opts.LogicalSize = runes.Measure(text)
	}

	if opts.Runes == "" {
		opts.Runes = marker.DefaultPaddingText
	}

	return opts
}

type LogicalSizeOpts struct {
	LogicalSize int
}

func fixLogicalSizeOpts(text string, opts LogicalSizeOpts) LogicalSizeOpts {
	if opts.LogicalSize == 0 {
		opts.LogicalSize = runes.Measure(text)
	}

	return opts
}

type TextTrimOpts struct {
	LogicalSize  int
	EllipsisText string
	EllipsisSize uint
}

func fixTextTrimOpts(text string, opts TextTrimOpts) TextTrimOpts {
	if opts.LogicalSize == 0 {
		opts.LogicalSize = runes.Measure(text)
	}

	if opts.EllipsisSize == 0 {
		opts.EllipsisSize = marker.DefaultElipsisSize
	}

	return opts
}

func Center(item any, width int) string {
	return CenterWithOpts(item, width, TextLayoutOpts{})
}

func CenterWithOpts(item any, width int, opts TextLayoutOpts) string {
	text := fmt.Sprintf("%v", item)

	opts = fixDirectionOps(text, opts)
	if opts.LogicalSize >= width {
		return text
	}

	padding := width - opts.LogicalSize
	left := padding / 2
	right := padding - left

	return strings.Repeat(opts.Runes, left) + text + strings.Repeat(opts.Runes, right)
}

func Left(item any, width int) string {
	return LeftWithOpts(item, width, TextLayoutOpts{})
}

func LeftWithOpts(item any, width int, opts TextLayoutOpts) string {
	text := fmt.Sprintf("%v", item)

	opts = fixDirectionOps(text, opts)
	if opts.LogicalSize >= width {
		return text
	}

	padding := width - opts.LogicalSize

	return strings.Repeat(opts.Runes, padding) + text
}

func Right(item any, width int) string {
	return RightWithOpts(item, width, TextLayoutOpts{})
}

func RightWithOpts(item any, width int, opts TextLayoutOpts) string {
	text := fmt.Sprintf("%v", item)

	opts = fixDirectionOps(text, opts)
	if opts.LogicalSize >= width {
		return text
	}

	padding := width - opts.LogicalSize

	return text + strings.Repeat(opts.Runes, padding)
}

func FillLeft(item any, width int) string {
	return FillLeftWithOpts(item, width, LogicalSizeOpts{})
}

func FillLeftWithOpts(item any, width int, opts LogicalSizeOpts) string {
	text := fmt.Sprintf("%v", item)

	opts = fixLogicalSizeOpts(text, opts)
	if opts.LogicalSize >= width {
		return text
	}

	if text == "" {
		text = marker.DefaultPaddingText
	}

	fix := ""
	if rest := width % len(text); rest != 0 {
		fix = text[rest:]
	}

	width = width / len(text)

	return fix + strings.Repeat(text, width)
}

func FillRight(item any, width int) string {
	return FillRightWithOpts(item, width, LogicalSizeOpts{})
}

func FillRightWithOpts(item any, width int, opts LogicalSizeOpts) string {
	text := fmt.Sprintf("%v", item)

	opts = fixLogicalSizeOpts(text, opts)
	if opts.LogicalSize >= width {
		return text
	}

	if text == "" {
		text = marker.DefaultPaddingText
	}

	fix := ""
	if rest := width % len(text); rest != 0 {
		fix = text[:rest]
	}

	width = width / len(text)

	return strings.Repeat(text, width) + fix
}

func RepeatLeft(item any, runes string, width int) string {
	return RepeatLeftWithOpts(item, runes, width, LogicalSizeOpts{})
}

func RepeatLeftWithOpts(item any, runes string, width int, opts LogicalSizeOpts) string {
	text := fmt.Sprintf("%v", item)
	return FillLeftWithOpts(runes, width, opts) + text
}

func RepeatRight(item any, runes string, width int) string {
	return RepeatRightWithOpts(item, runes, width, LogicalSizeOpts{})
}

func RepeatRightWithOpts(item any, runes string, width int, opts LogicalSizeOpts) string {
	text := fmt.Sprintf("%v", item)
	return text + FillRightWithOpts(runes, width, opts)
}

func TrimLeft(data string, width int, opts TextTrimOpts) string {
	if data == "" {
		return data
	}

	opts = fixTextTrimOpts(data, opts)

	elipSize := runes.Measure(opts.EllipsisText) * int(opts.EllipsisSize)

	realSize := runes.Measure(data)
	if width >= opts.LogicalSize || width > realSize {
		return data
	}

	width = realSize - width

	if elipSize+width >= realSize {
		index, _ := runes.RuneIndexToByteIndex(data, width)
		return data[index:]
	}

	elipTotal := strings.Repeat(opts.EllipsisText, int(opts.EllipsisSize))

	index, _ := runes.RuneIndexToByteIndex(data, width+elipSize)
	return elipTotal + data[index:]
}

func TrimRight(data string, width int, opts TextTrimOpts) string {
	if data == "" {
		return data
	}

	opts = fixTextTrimOpts(data, opts)

	elipSize := runes.Measure(opts.EllipsisText) * int(opts.EllipsisSize)

	realSize := runes.Measure(data)
	if width >= opts.LogicalSize || width > realSize {
		return data
	}

	if elipSize > width {
		index, _ := runes.RuneIndexToByteIndex(data, width)
		return data[:index]
	}

	elipTotal := strings.Repeat(opts.EllipsisText, int(opts.EllipsisSize))

	index, _ := runes.RuneIndexToByteIndex(data, width-elipSize)
	return data[:index] + elipTotal
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
