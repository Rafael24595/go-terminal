package input

type CheckActionHandler = func()

func voidCheckHandler() {}

type CheckAction struct {
	ActionMode bool
	Handler    CheckActionHandler
}

func NewCheckAction(handler ...CheckActionHandler) *CheckAction {
	hdl := voidCheckHandler
	if len(handler) > 0 {
		hdl = handler[0]
	}

	return &CheckAction{
		ActionMode: false,
		Handler:    hdl,
	}
}
