package rpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRpcCall(t *testing.T) {
	adapter := NewMemoryAdapter()

	done, err := StartRpcGreetingServiceServer(&RpcGreetingServiceImpl{}, adapter)
	assert.NoError(t, err)
	defer done()

	client := CreateRpcGreetingServiceClient(adapter)

	resp, err := client.SayHello(context.Background(), "World")
	assert.NoError(t, err)
	assert.Equal(t, "Hello World", resp)
}
