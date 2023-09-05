package http

import (
	"context"
	"sync"
)

type HttpSignalServiceImpl struct {
	bus *MemoryBus
}

func (h *HttpSignalServiceImpl) Send(ctx context.Context, inbox string, msg string) (err error) {
	return h.bus.Send(ctx, inbox, msg)
}

func (h *HttpSignalServiceImpl) Recv(ctx context.Context, inbox string) (msgs <-chan string, err error) {
	return h.bus.Recv(ctx, inbox)
}

func NewHttpSignalServiceImpl() *HttpSignalServiceImpl {
	return &HttpSignalServiceImpl{
		bus: &MemoryBus{
			inboxes: make(map[string]chan string),
		},
	}
}

type MemoryBus struct {
	inboxes map[string]chan string
	mux     sync.Mutex
}

func (m *MemoryBus) Send(ctx context.Context, inbox string, msg string) (err error) {
	inboxChannel := m.getInboxChannel(inbox)
	inboxChannel <- msg
	return nil
}

func (m *MemoryBus) Recv(ctx context.Context, inbox string) (msgs <-chan string, err error) {
	inboxChannel := m.getInboxChannel(inbox)
	go func() {
		<-ctx.Done()
		close(inboxChannel)
	}()

	return inboxChannel, nil
}

func (m *MemoryBus) getInboxChannel(inbox string) chan string {
	m.mux.Lock()
	defer m.mux.Unlock()

	inboxChannel, ok := m.inboxes[inbox]
	if !ok {
		inboxChannel = make(chan string, 10)
		m.inboxes[inbox] = inboxChannel
	}

	return inboxChannel
}
