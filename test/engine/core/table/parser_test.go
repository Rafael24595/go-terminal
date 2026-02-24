package table_test

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core/table"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

type testStruct struct {
	ID      int
	Name    string
	Project bool
}

type emptyStruct struct{}

type pointerStruct struct {
	Lang string
}

func TestStructFieds_ShouldReturnAllFields(t *testing.T) {
	s := testStruct{
		ID:      10,
		Name:    "Golang",
		Project: true,
	}

	fields := table.StructFieds(s)

	assert.Equal(t, 3, len(fields))

	assert.Equal(t, "ID", fields[0].Header)
	assert.Equal(t, 10, fields[0].Value)

	assert.Equal(t, "Name", fields[1].Header)
	assert.Equal(t, "Golang", fields[1].Value)

	assert.Equal(t, "Project", fields[2].Header)
	assert.Equal(t, true, fields[2].Value)
}

func TestStructFieds_ShouldWorkWithPointer(t *testing.T) {
	s := &pointerStruct{
		Lang: "ziglang",
	}

	fields := table.StructFieds(s)

	assert.Equal(t, 1, len(fields))
	assert.Equal(t, "Lang", fields[0].Header)
	assert.Equal(t, "ziglang", fields[0].Value)
}

func TestStructFieds_EmptyStruct_ShouldReturnEmptySlice(t *testing.T) {
	s := emptyStruct{}

	fields := table.StructFieds(s)

	assert.Equal(t, 0, len(fields))
}

func TestStructHeaders_ShouldReturnOnlyHeaders(t *testing.T) {
	headers := table.StructHeaders[testStruct]()

	assert.Equal(t, 3, len(headers))

	assert.Equal(t, "ID", headers[0])
	assert.Equal(t, "Name", headers[1])
	assert.Equal(t, "Project", headers[2])
}

func TestStructFieds_Nil_ShouldReturnEmpty(t *testing.T) {
	var s *pointerStruct

	assert.NotPanic(t, func() {
		table.StructFieds(s)
	})
}

func TestStructFieds_NonStruct_ShouldReturnNil(t *testing.T) {
	fields := table.StructFieds(123)

	assert.Len(t, 0, fields)
}
