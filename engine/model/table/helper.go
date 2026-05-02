package table

func Rows(headers []string, cols map[string][]string) uint16 {
	colSize := uint16(0)
	for _, h := range headers {
		colSize = max(
			colSize, uint16(len(cols[h])),
		)
	}
	return colSize
}
