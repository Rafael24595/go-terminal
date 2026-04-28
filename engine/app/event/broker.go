package event

import (
	"context"
	"sync"

	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
)

const default_buffer_size = uint(16)

type subscriber[T any] struct {
	channel   chan Event[T]
	timestamp int64
	status    bool
}

type EventBroker[T any] struct {
	clock clock.Clock
	init  sync.Once
	mutx  sync.RWMutex
	chnl  chan Event[T]
	done  chan struct{}
	subs  *dict.LinkedMap[chan Event[T], subscriber[T]]
}

func NewBroker[T any]() *EventBroker[T] {
	return new(EventBroker[T])
}

func NewBrokerWithCtx[T any](ctx context.Context, buffer ...uint) *EventBroker[T] {
	return new(EventBroker[T]).Init(ctx, buffer...)
}

func (b *EventBroker[T]) lazyInit() *EventBroker[T] {
	return b.Init(context.Background())
}

func (b *EventBroker[T]) Init(ctx context.Context, buffer ...uint) *EventBroker[T] {
	b.init.Do(func() {
		bufferSize := default_buffer_size
		if len(buffer) > 0 {
			bufferSize = buffer[0]
		}

		b.chnl = make(chan Event[T], bufferSize)
		b.done = make(chan struct{})
		b.clock = clock.UnixMilliClock
		b.subs = dict.NewLinkedMap[chan Event[T], subscriber[T]]()

		go b.dispatch()
		go b.listenContext(ctx)
	})
	return b
}

func (b *EventBroker[T]) Close() {
	b.lazyInit()

	b.mutx.Lock()
	defer b.mutx.Unlock()

	select {
	case <-b.done:
	default:
		close(b.done)
	}
}

func (b *EventBroker[T]) Subscribe(subs ...chan Event[T]) {
	b.lazyInit()

	if b.isDone() {
		assert.Unreachable("no new elements can subscribed to closed broker")
		return
	}

	b.mutx.Lock()
	defer b.mutx.Unlock()

	for _, s := range subs {
		if b.subs.Exists(s) {
			continue
		}

		b.subs.Set(s, subscriber[T]{
			channel:   s,
			timestamp: b.clock(),
			status:    true,
		})
	}
}

func (b *EventBroker[T]) Unsubscribe(subs ...chan Event[T]) {
	b.lazyInit()

	if b.isDone() {
		assert.Unreachable("cannot unsubscribe from a closed broker")
		return
	}

	b.mutx.Lock()
	defer b.mutx.Unlock()

	for _, s := range subs {
		b.subs.Delete(s)
	}
}

func (b *EventBroker[T]) Publish(value T) {
	b.lazyInit()

	select {
	case b.chnl <- NewEvent(value):
		return
	case <-b.done:
		assert.Unreachable("cannot publish to a closed broker")
	}
}

func (b *EventBroker[T]) IsRunning() bool {
	if b.done == nil {
		return false
	}
	return !b.isDone()
}

func (b *EventBroker[T]) isDone() bool {
	select {
	case <-b.done:
		return true
	default:
		return false
	}
}

func (b *EventBroker[T]) dispatch() {
	for {
		select {
		case e, ok := <-b.chnl:
			if !ok {
				return
			}
			b.broadcast(e)
		case <-b.done:
			return
		}
	}
}

func (b *EventBroker[T]) broadcast(e Event[T]) {
	b.mutx.RLock()
	defer b.mutx.RUnlock()

	for s := range b.subs.Values() {
		select {
		case s.channel <- e:
		default:
		}
	}
}

func (b *EventBroker[T]) listenContext(ctx context.Context) {
	select {
	case <-ctx.Done():
		b.Close()
	case <-b.done:
		return
	}
}
