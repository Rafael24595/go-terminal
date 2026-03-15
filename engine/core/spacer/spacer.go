package spacer

type SpacerMode uint8

const (
	SpacerAppend SpacerMode = iota
	SpacerAfterEach
)

type SpacerMeta struct {
	Size uint8
	Mode SpacerMode
}

func NewSpacerMeta(size uint8, mode SpacerMode) SpacerMeta {
	return SpacerMeta{
		Size: size,
		Mode: mode,
	}
}
