package main

import (
	"context"
	"net/http"

	"ella.to/examples/http/schema"
)

func main() {
	host := "http://localhost:8080"
	httpClient := &http.Client{}

	client := schema.CreateHttpGreetingServiceClient(host, httpClient)

	result, err := client.Hello(context.Background(), "Ella")
	if err != nil {
		panic(err)
	}

	println(result)
}
