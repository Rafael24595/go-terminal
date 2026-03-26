package input

type TableActionHandler = func(MatrixCursor)

func voidTableHandler(_ MatrixCursor) {}

type TableAction struct {
	EnableMode bool
	ActionMode bool
	Handler    TableActionHandler
}

func NewTableAction(handler ...TableActionHandler) *TableAction {
	enb := false
	hdl := voidTableHandler
	if len(handler) > 0 {
		enb = true
		hdl = handler[0]
	}

	return &TableAction{
		EnableMode: enb,
		ActionMode: false,
		Handler:    hdl,
	}
}

func (a *TableAction) SetHandler(handler TableActionHandler) *TableAction {
	a.EnableMode = true
	a.Handler = handler
	return a
}
