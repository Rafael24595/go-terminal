package runes_test

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/helper/runes"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestAppendAt(t *testing.T) {
	tests := []struct {
		name   string
		slice  []rune
		insert []rune
		pos    uint
		want   string
	}{
		{
			name:   "insert in the middle",
			slice:  []rune("hello world"),
			insert: []rune(" beautiful"),
			pos:    5,
			want:   "hello beautiful world",
		},
		{
			name:   "insert at beginning",
			slice:  []rune("world"),
			insert: []rune("hello "),
			pos:    0,
			want:   "hello world",
		},
		{
			name:   "insert at end",
			slice:  []rune("hello"),
			insert: []rune(" world"),
			pos:    5,
			want:   "hello world",
		},
		{
			name:   "insert empty slice",
			slice:  []rune("hello"),
			insert: []rune(""),
			pos:    3,
			want:   "hello",
		},
		{
			name:   "insert into empty slice",
			slice:  []rune(""),
			insert: []rune("hello"),
			pos:    0,
			want:   "hello",
		},
		{
			name:   "unicode runes",
			slice:  []rune("hola 🌍"),
			insert: []rune(" querido"),
			pos:    4,
			want:   "hola querido 🌍",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runes.AppendAt(tt.slice, tt.insert, tt.pos)

			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestAppendRange(t *testing.T) {
	tests := []struct {
		name   string
		slice  []rune
		insert []rune
		start  uint
		end    uint
		want   string
	}{
		{
			name:   "replace range in the middle",
			slice:  []rune("hello golang"),
			insert: []rune("zig"),
			start:  6,
			end:    8,
			want:   "hello ziglang",
		},
		{
			name:   "replace at beginning",
			slice:  []rune("hello world"),
			insert: []rune("hey"),
			start:  0,
			end:    5,
			want:   "hey world",
		},
		{
			name:   "replace at end",
			slice:  []rune("hello world"),
			insert: []rune("gophers"),
			start:  6,
			end:    11,
			want:   "hello gophers",
		},
		{
			name:   "insert when start equals end",
			slice:  []rune("hello world"),
			insert: []rune(" beautiful"),
			start:  5,
			end:    5,
			want:   "hello beautiful world",
		},
		{
			name:   "delete range (empty insert)",
			slice:  []rune("hello cruel world"),
			insert: []rune(""),
			start:  6,
			end:    12,
			want:   "hello world",
		},
		{
			name:   "replace with longer text",
			slice:  []rune("hello go"),
			insert: []rune("golang language"),
			start:  6,
			end:    8,
			want:   "hello golang language",
		},
		{
			name:   "replace with shorter text",
			slice:  []rune("hello golang language"),
			insert: []rune("go"),
			start:  6,
			end:    21,
			want:   "hello go",
		},
		{
			name:   "replace entire slice",
			slice:  []rune("hello world"),
			insert: []rune("goodbye"),
			start:  0,
			end:    11,
			want:   "goodbye",
		},
		{
			name:   "unicode runes",
			slice:  []rune("hi 🐿️"),
			insert: []rune(" dear"),
			start:  2,
			end:    3,
			want:   "hi dear🐿️",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runes.AppendRange(tt.slice, tt.insert, tt.start, tt.end)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

