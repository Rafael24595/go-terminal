package event

import (
	"context"
	"sync"
	"testing"
	"time"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestBroker_ZeroValueUsage(t *testing.T) {
	var b EventBroker[int]

	ch := make(chan Event[int], 1)

	b.Subscribe(ch)
	b.Publish(42)

	select {
	case e := <-ch:
		assert.Equal(t, 42, e.Value)
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout: zero value did not initialize the goroutines correctly")
	}
}

func TestEventBroker_PublishSubscribe(t *testing.T) {
	broker := NewBroker[string]()
	ch := make(chan Event[string], 1)

	broker.Subscribe(ch)

	expected := "golang"
	broker.Publish(expected)

	select {
	case e := <-ch:
		assert.Equal(t, expected, e.Value)
	case <-time.After(time.Second):
		t.Fatal("timeout: the event never reached the subscriber")
	}
}

func TestEventBroker_MultipleSubscribers(t *testing.T) {
	broker := NewBroker[int]()

	const numSubs = 10
	const numEvents = 5

	var wg sync.WaitGroup
	wg.Add(numSubs * numEvents)

	channels := make([]chan Event[int], numSubs)
	for i := range numSubs {
		channels[i] = make(chan Event[int], numEvents)
		broker.Subscribe(channels[i])

		go func(c chan Event[int]) {
			for range c {
				wg.Done()
			}
		}(channels[i])
	}

	for i := range numEvents {
		broker.Publish(i)
	}

	c := make(chan struct{})
	go func() {
		wg.Wait()
		c <- struct{}{}
	}()

	select {
	case <-c:
	case <-time.After(time.Second):
		t.Fatal("timeout: waiting for subscribers")
	}
}

func TestEventBroker_Unsubscribe(t *testing.T) {
	broker := NewBroker[int]()
	ch := make(chan Event[int], 1)

	broker.Subscribe(ch)
	broker.Unsubscribe(ch)

	broker.Publish(100)

	select {
	case <-ch:
		t.Fatal("the subscriber received a message after unsubscribing.")
	case <-time.After(time.Millisecond * 50):
	}
}

func TestEventBroker_SlowSubscriberDoesNotBlock(t *testing.T) {
	broker := NewBroker[string]()

	chSlow := make(chan Event[string])
	chFast := make(chan Event[string], 1)

	broker.Subscribe(chFast, chSlow)

	done := make(chan bool)
	go func() {
		broker.Publish("ziglang")
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(time.Millisecond * 100):
		t.Fatal("timeout: the broker was blocked due to the slow subscriber")
	}

	select {
	case <-chFast:
	case <-time.After(time.Second):
		t.Fatal("timeout: the event never reached the subscriber")
	}
}

func TestEventBroker_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	broker := NewBrokerWithCtx[int](ctx)

	cancel()

	<-broker.done

	assert.False(t, broker.IsRunning())
}

func TestBroker_MultipleCloseSources(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	broker := NewBrokerWithCtx[int](ctx)

	go cancel()
	go broker.Close()

	<-broker.done

	assert.False(t, broker.IsRunning())

	assert.Panic(t, func() {
		broker.Publish(10)
	})
}

func TestBroker_SubscribeUnsubscribeStress(t *testing.T) {
	broker := NewBroker[int]()
	channels := make([]chan Event[int], 100)

	for i := range channels {
		channels[i] = make(chan Event[int], 1)
		broker.Subscribe(channels[i])
	}

	for i := range 50 {
		broker.Unsubscribe(channels[i])
	}

	broker.Publish(1)

	for i := range 50 {
		select {
		case <-channels[i]:
			t.Errorf("the %d channel received an event", i)
		default:
		}
	}
}
