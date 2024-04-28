regenrate:
	go run main.go gen http ./e2e/http/http.gen.go ./e2e/http/http.ella
	go run main.go gen stream ./e2e/stream/stream.gen.go ./e2e/stream/stream.ella
	go run main.go gen upload ./e2e/upload/upload.gen.go ./e2e/upload/upload.ella
	go run main.go gen rpc ./e2e/rpc/rpc.gen.go ./e2e/rpc/rpc.ella
	go run main.go gen http ./e2e/http_async_stream/http_async_stream.gen.go ./e2e/http_async_stream/http_async_stream.ella
	go run main.go gen download ./e2e/download/download.gen.go ./e2e/download/download.ella

run-e2e: regenrate
	go test ./e2e/http/... -v
	go test ./e2e/stream/... -v
	go test ./e2e/upload/... -v
	go test ./e2e/rpc/... -v
	go test ./e2e/http_async_stream/... -v
	go test ./e2e/download/... -v