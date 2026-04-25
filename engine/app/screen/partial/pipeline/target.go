package pipeline

type Target uint8

const (
	Header Target = iota
	Kernel
	Footer
)
