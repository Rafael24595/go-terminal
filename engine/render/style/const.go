package style

type Direction uint8

const (
	Horizontal Direction = iota
	Vertical
)

type VerticalPosition uint

const (
	Left VerticalPosition = iota
	Center
	Right
)

type HorizontalPosition uint

const (
	Top HorizontalPosition = iota
	Middle
	Bottom
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

