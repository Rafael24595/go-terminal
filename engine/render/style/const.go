package style

type Direction uint8

const (
	Horizontal Direction = iota
	Vertical
)

type VerticalPosition uint

const (
	Top VerticalPosition = iota
	Middle
	Bottom
)

type HorizontalPosition uint

const (
	Left HorizontalPosition = iota
	Center
	Right
)

type Justify uint8

const (
	JustifyStart Justify = iota
	JustifyEnd
	JustifyCenter
	JustifyBetween
	JustifyAround
	JustifyEvenly
)
