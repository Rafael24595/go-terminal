package screen

import "github.com/Rafael24595/go-reacterm-core/engine/model/key"

type ScreenEvent struct {
	Key key.Key
}

func NewEvent(key key.Key) ScreenEvent {
	return ScreenEvent{
		Key: key,
	}
}
