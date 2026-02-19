package table

func Cols(headers []string, cols map[string][]string) int {
	colSize := 0
	for _, h := range headers {
		if len(cols[h]) > colSize {
			colSize = len(cols[h])
		}
	}
	return colSize
}
