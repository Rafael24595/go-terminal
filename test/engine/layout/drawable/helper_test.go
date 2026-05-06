package drawable_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func TestDrainDrawable_WithMock(t *testing.T) {
	tests := []struct {
		name          string
		rows          winsize.Rows
		batch         uint
		lines         int
		lazy          bool
		wantLines     int
		wantDrawCalls int
	}{
		{
			name:          "Eager_DrainAllInSingleCall",
			rows:          10,
			batch:         0,
			lines:         5,
			lazy:          false,
			wantLines:     5,
			wantDrawCalls: 1,
		},
		{
			name:          "Eager_DrainEverythingAcrossMultipleCalls",
			rows:          2,
			batch:         2,
			lines:         6,
			lazy:          false,
			wantLines:     6,
			wantDrawCalls: 3,
		},
		{
			name:          "Lazy_StopExactlyAtRowsLimit",
			rows:          4,
			batch:         2,
			lines:         10,
			lazy:          true,
			wantLines:     4,
			wantDrawCalls: 2,
		},
		{
			name:          "Lazy_StopWhenDrawableIsExhaustedBeforeLimit",
			rows:          10,
			batch:         2,
			lines:         4,
			lazy:          true,
			wantLines:     4,
			wantDrawCalls: 2,
		},
		{
			name:          "EdgeCase_ZeroRowsWithLazyMode",
			rows:          0,
			batch:         5,
			lines:         10,
			lazy:          true,
			wantLines:     0,
			wantDrawCalls: 1,
		},
		{
			name:          "EdgeCase_EmptyDrawable",
			rows:          10,
			batch:         5,
			lines:         0,
			lazy:          false,
			wantLines:     0,
			wantDrawCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDrawable{
				Lines: make([]text.Line, tt.lines),
				Batch: tt.batch,
			}

			size := winsize.Winsize{
				Rows: tt.rows,
			}

			got, _ := drawable.DrainDrawable(size, m.ToDrawable(), tt.lazy)

			assert.Len(t, tt.wantLines, got)
			assert.Equal(t, tt.wantDrawCalls, m.DrawCalls)
		})
	}
}
