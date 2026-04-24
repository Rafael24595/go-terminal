package chunk

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func TestDynamic(t *testing.T) {
	c := Dynamic()

	assert.False(t, c.Sized)

	res := c.Adapter(winsize.Winsize{
		Cols: 80,
	})

	assert.Equal(t, 0, res)
}

func TestColums(t *testing.T) {
	tests := []struct {
		name     string
		columns  uint16
		terminal uint16
		expected uint16
	}{
		{"Normal", 20, 80, 20},
		{"Exact fit", 80, 80, 80},
		{"Clamping (Overflow)", 100, 80, 80},
		{"Zero", 0, 80, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Colums(tt.columns)
			assert.True(t, c.Sized)

			res := c.Adapter(winsize.Winsize{
				Cols: tt.terminal,
			})

			assert.Equal(t, tt.expected, res)
		})
	}
}

func TestPercent(t *testing.T) {
	tests := []struct {
		name     string
		percent  uint16
		terminal uint16
		expected uint16
	}{
		{"Half", 50, 100, 50},
		{"Quarter", 25, 80, 20},
		{"Full", 100, 120, 120},
		{"Zero", 0, 80, 0},
		{"Rounding", 33, 100, 33},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Percent(tt.percent)
			assert.True(t, c.Sized)

			res := c.Adapter(winsize.Winsize{Cols: tt.terminal})

			assert.Equal(t, tt.expected, res)
		})
	}
}

func TestPercent_Panic(t *testing.T) {
	assert.Panic(t, func() {
		Percent(101)
	})
}
