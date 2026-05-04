package pipeline

type Section uint8

const (
	Header Section = iota
	Kernel
	Footer
)
