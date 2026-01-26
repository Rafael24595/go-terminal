package core

type Screen struct {
	//Init func (ctx)
	//Update func(evt)
	View func() ViewModel
}