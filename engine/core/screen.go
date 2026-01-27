package core

type ScreenEvent struct {
	Key string
}

type Screen struct {
	//Init func (ctx)
	Update func(ScreenEvent)
	View   func() ViewModel
}
