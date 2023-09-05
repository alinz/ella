package stream

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpStream(t *testing.T) {
	server := httptest.NewServer(
		CreateEventServiceServer(&HttpEventServiceImpl{}),
	)
	defer server.Close()

	host := server.URL
	httpClient := &http.Client{}

	client := CreateHttpEventServiceClient(host, httpClient)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	results, err := client.GetRandomValues(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for result := range results {
		println(result)
	}
}
