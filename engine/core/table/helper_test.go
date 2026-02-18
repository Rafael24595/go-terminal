package table

import (
	"testing"

	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestAdjustSize_NoReductionNeeded(t *testing.T) {
	size := map[string]int{
		"A": 5,
		"B": 5,
	}

	separator := SeparatorMeta{
		center: "|",
		left:   "|",
		right:  "|",
	}

	termWidth := 20

	rendered := renderedRowSize(size, separator)
	result, status := adjustSize(size, termWidth, rendered)

	assert.True(t, status)

	assert.Equal(t, size["A"], result["A"])
	assert.Equal(t, size["B"], result["B"])
}

func TestAdjustSize_ReducesLargestColumn(t *testing.T) {
	size := map[string]int{
		"A": 10,
		"B": 5,
	}

	separator := SeparatorMeta{
		center: "|",
		left:   "|",
		right:  "|",
	}

	termWidth := 14

	rendered := renderedRowSize(size, separator)
	result, status := adjustSize(size, termWidth, rendered)

	assert.True(t, status)

	assert.GreaterOrEqual(t, termWidth, renderedRowSize(result, separator))
	assert.Less(t, 10, result["A"])
}

func TestAdjustSize_RespectsMinWidth(t *testing.T) {
	size := map[string]int{
		"A": 4,
		"B": 4,
	}

	separator := SeparatorMeta{
		center: "|",
		left:   "|",
		right:  "|",
	}
	termWidth := 5

	rendered := renderedRowSize(size, separator)
	result, status := adjustSize(size, termWidth, rendered)

	assert.False(t, status)

	assert.GreaterOrEqual(t, 4, result["A"])
	assert.GreaterOrEqual(t, 4, result["B"])
}

func TestAdjustSize_ExactFit(t *testing.T) {
	size := map[string]int{
		"A": 8,
		"B": 6,
	}

	separator := SeparatorMeta{
		center: "|",
		left:   "|",
		right:  "|",
	}

	rendered := renderedRowSize(size, separator)

	termWidth := rendered - 3

	result, status := adjustSize(size, termWidth, rendered)

	assert.True(t, status)

	assert.Equal(t, termWidth, renderedRowSize(result, separator))
}

func TestAdjustSize_MultipleColumnsReduction(t *testing.T) {
	size := map[string]int{
		"A": 10,
		"B": 9,
		"C": 8,
	}

	separator := SeparatorMeta{
		center: "|",
		left:   "|",
		right:  "|",
	}

	termWidth := 20

	rendered := renderedRowSize(size, separator)
	result, status := adjustSize(size, termWidth, rendered)

	assert.True(t, status)

	assert.Equal(t, termWidth, renderedRowSize(result, separator))

	assert.NotEqual(t, 10, result["A"])
	assert.NotEqual(t, 9, result["B"])
	assert.NotEqual(t, 8, result["C"])
}

func TestSplitTable_FitsInOne(t *testing.T) {
	size := map[string]int{
		"A": 10,
		"B": 20,
		"C": 10,
	}

	termWidth := 50

	result := splitTable(size, termWidth)

	assert.Equal(t, 1, len(result))
	assert.Equal(t, 10, result[0]["A"])
	assert.Equal(t, 20, result[0]["B"])
	assert.Equal(t, 10, result[0]["C"])
}

func TestSplitTable_MustSplit(t *testing.T) {
	size := map[string]int{
		"A": 20,
		"B": 10,
		"C": 15,
		"D": 15,
	}

	termWidth := 25

	result := splitTable(size, termWidth)

	assert.True(t, len(result) > 1)

	for _, table := range result {
		total := 0
		for _, v := range table {
			total += v
		}
		assert.True(t, total <= termWidth)
	}
}

func TestSplitTable_ColumnWiderThanTerminal(t *testing.T) {
	size := map[string]int{
		"XL": 100,
	}

	termWidth := 80

	result := splitTable(size, termWidth)

	assert.Equal(t, 1, len(result))
	assert.Equal(t, 80, result[0]["XL"])
}

func TestSplitTable_EmptyMap(t *testing.T) {
	size := map[string]int{}
	result := splitTable(size, 80)
	assert.Equal(t, 0, len(result))
}
