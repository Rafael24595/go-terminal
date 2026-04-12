package help

type HelpField struct {
	Code   []string
	Detail string
}

type HelpMeta struct {
	Show   bool
	Fields []HelpField
}

func NewHelpMeta() *HelpMeta {
	return &HelpMeta{
		Show:   false,
		Fields: make([]HelpField, 0),
	}
}

func (d *HelpMeta) Unshift(fields ...HelpField) *HelpMeta {
	d.Fields = append(fields, d.Fields...)
	return d
}

func (d *HelpMeta) Push(fields ...HelpField) *HelpMeta {
	d.Fields = append(d.Fields, fields...)
	return d
}
