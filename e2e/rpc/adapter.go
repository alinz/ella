package rpc

import (
	"context"
	"sync"
)

type AdapterMsg struct {
	topic   string
	data    []byte
	replyFn func([]byte) error
}

var _ rpcMsg = (*AdapterMsg)(nil)

func (a *AdapterMsg) Topic() string {
	return a.topic
}

func (a *AdapterMsg) Data() []byte {
	return a.data
}

func (a *AdapterMsg) Reply(data []byte) error {
	return a.replyFn(data)
}

type Adapter struct {
	mux    sync.Mutex
	topics map[string]chan *AdapterMsg
}

var _ rpcAdaptor = (*Adapter)(nil)

func (a *Adapter) Register(topic string, recv recvFunc) (drain func(), err error) {
	topicChannel := a.getTopicChannel(topic)

	go func() {
		for msg := range topicChannel {
			recv(msg)
		}
	}()

	return func() {
		a.mux.Lock()
		defer a.mux.Unlock()

		close(topicChannel)
		delete(a.topics, topic)
	}, nil
}

func (a *Adapter) Send(ctx context.Context, topic string, data []byte) ([]byte, error) {
	resp := make(chan *AdapterMsg, 1)

	topicChannel := a.getTopicChannel(topic)
	topicChannel <- &AdapterMsg{
		topic: topic,
		data:  data,
		replyFn: func(data []byte) error {
			resp <- &AdapterMsg{
				topic: topic,
				data:  data,
			}
			return nil
		},
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg := <-resp:
		return msg.data, nil
	}
}

func (a *Adapter) getTopicChannel(name string) chan *AdapterMsg {
	a.mux.Lock()
	defer a.mux.Unlock()

	if ch, ok := a.topics[name]; ok {
		return ch
	}

	ch := make(chan *AdapterMsg, 10)
	a.topics[name] = ch
	return ch
}

func NewMemoryAdapter() *Adapter {
	return &Adapter{
		topics: make(map[string]chan *AdapterMsg),
	}
}
