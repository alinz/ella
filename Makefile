regenrate:
	go run cmd/ella/main.go gen http ./e2e/http/http.gen.go ./e2e/http/http.ella
	go run cmd/ella/main.go gen stream ./e2e/stream/stream.gen.go ./e2e/stream/stream.ella
	go run cmd/ella/main.go gen upload ./e2e/upload/upload.gen.go ./e2e/upload/upload.ella

run-e2e: regenrate
	go test ./e2e/http/... -v
	go test ./e2e/stream/... -v
	go test ./e2e/upload/... -v