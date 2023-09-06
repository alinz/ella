package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStreamWithAsync(t *testing.T) {
	server := httptest.NewServer(
		CreateSignalServiceServer(
			NewHttpSignalServiceImpl(
				NewMemoryBus[string](),
			),
		),
	)
	defer server.Close()

	host := server.URL

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		client := CreateHttpSignalServiceClient(host, &http.Client{})

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		err := client.Send(ctx, "inbox", "Hello")
		assert.NoError(t, err)
	}()

	go func() {
		httpClient := &http.Client{}
		defer func() {
			wg.Done()
		}()

		client := CreateHttpSignalServiceClient(host, httpClient)

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		msgs, err := client.Recv(ctx, "inbox")
		assert.NoError(t, err)

		msg := <-msgs
		assert.Equal(t, "Hello", msg)
	}()

	wg.Wait()
}
