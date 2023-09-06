package http

import (
	"context"
	"sync"
)

type Bus[T any] interface {
	Send(ctx context.Context, inbox string, msg T) (err error)
	Recv(ctx context.Context, inbox string) (msgs <-chan T, err error)
}

type HttpSignalServiceImpl struct {
	bus Bus[string]
}

func (h *HttpSignalServiceImpl) Send(ctx context.Context, inbox string, msg string) (err error) {
	return h.bus.Send(ctx, inbox, msg)
}

func (h *HttpSignalServiceImpl) Recv(ctx context.Context, inbox string) (msgs <-chan string, err error) {
	return h.bus.Recv(ctx, inbox)
}

func NewHttpSignalServiceImpl(bus Bus[string]) *HttpSignalServiceImpl {
	return &HttpSignalServiceImpl{
		bus: bus,
	}
}

type MemoryBus[T any] struct {
	inboxes map[string]chan T
	mux     sync.Mutex
}

var _ Bus[string] = (*MemoryBus[string])(nil)

func (m *MemoryBus[T]) Send(ctx context.Context, inbox string, msg T) (err error) {
	inboxChannel := m.getInboxChannel(inbox)
	inboxChannel <- msg
	return nil
}

func (m *MemoryBus[T]) Recv(ctx context.Context, inbox string) (msgs <-chan T, err error) {
	inboxChannel := m.getInboxChannel(inbox)
	go func() {
		<-ctx.Done()
		close(inboxChannel)
	}()

	return inboxChannel, nil
}

func (m *MemoryBus[T]) getInboxChannel(inbox string) chan T {
	m.mux.Lock()
	defer m.mux.Unlock()

	inboxChannel, ok := m.inboxes[inbox]
	if !ok {
		inboxChannel = make(chan T, 10)
		m.inboxes[inbox] = inboxChannel
	}

	return inboxChannel
}

func NewMemoryBus[T any]() *MemoryBus[T] {
	return &MemoryBus[T]{
		inboxes: map[string]chan T{},
	}
}
