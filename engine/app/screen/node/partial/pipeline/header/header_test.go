package header

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func TestHeader_InsertsBefore(t *testing.T) {
	vm := viewmodel.ViewModel{
		Header: stack.NewVStack(
			line.DrawableFromLines(
				*text.NewLine("line_01"),
			),
		),
	}

	items := vm.Header.Items()
	assert.Len(t, 1, items)

	line := text.NewLine("line_02")
	transformer := Transformer(pipeline.Before, *line)
	vm = transformer(vm)

	items = vm.Header.Items()
	assert.Len(t, 2, items)

	drw := items[0]
	drw.Init()
	lns, _ := drw.Draw(winsize.Winsize{
		Rows: 1,
		Cols: 10,
	})

	assert.Equal(t, "line_02", text.LineToString(&lns[0]))
}

func TestHeader_InsertsAfter(t *testing.T) {
	vm := viewmodel.ViewModel{
		Header: stack.NewVStack(
			line.DrawableFromLines(
				*text.NewLine("line_01"),
			),
		),
	}

	items := vm.Header.Items()
	assert.Len(t, 1, items)

	line := text.NewLine("line_02")
	transformer := Transformer(pipeline.After, *line)
	vm = transformer(vm)

	items = vm.Header.Items()
	assert.Len(t, 2, items)

	drw := items[1]
	drw.Init()
	lns, _ := drw.Draw(winsize.Winsize{
		Rows: 1,
		Cols: 10,
	})

	assert.Equal(t, "line_02", text.LineToString(&lns[0]))
}
