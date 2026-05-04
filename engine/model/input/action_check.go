package input

type CheckActionHandler = func()

func voidCheckHandler() {}

type CheckAction struct {
	ActionMode bool
	Handler    CheckActionHandler
}

func NewCheckAction(handler CheckActionHandler) *CheckAction {
	return &CheckAction{
		ActionMode: false,
		Handler:    handler,
	}
}

func EmptyCheckAction() *CheckAction {
	return NewCheckAction(voidCheckHandler)
}
