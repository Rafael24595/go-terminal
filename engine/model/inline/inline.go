package inline

const DefaultInlineSeparator = " | "

type Target uint8

const (
	TargetCode Target = iota
	TargetTags
)

type FilterMeta struct {
	Target Target
	Values []string
}

func NewFilterMeta(target Target, values ...string) FilterMeta {
	return FilterMeta{
		Target: target,
		Values: values,
	}
}
