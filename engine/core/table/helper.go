package table

func Rows(headers []string, cols map[string][]string) int {
	colSize := 0
	for _, h := range headers {
		colSize = max(colSize, len(cols[h]))
	}
	return colSize
}
