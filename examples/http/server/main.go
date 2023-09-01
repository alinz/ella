package main

import (
	"context"
	"fmt"
	"net/http"

	"ella.to/examples/http/schema"
)

type httpGreetingService struct{}

var _ schema.HttpGreetingService = (*httpGreetingService)(nil)

func (s *httpGreetingService) Hello(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("Hello %s", name), nil
}

func main() {
	serverHandler := schema.CreateGreetingServiceServer(&httpGreetingService{})
	http.Handle("/", serverHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
