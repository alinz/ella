package download_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"compiler.ella.to/e2e/download"
)

func TestCallHttpMethod(t *testing.T) {
	server := httptest.NewServer(
		download.CreateDownloadServiceServer(&download.HttpDownloadServiceImpl{}),
	)

	host := server.URL
	fmt.Println(host)

	httpClient := &http.Client{}

	client := download.CreateHttpDownloadServiceClient(host, httpClient)

	r, filename, contentType, err := client.Get(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "text/plain", contentType)
	assert.Equal(t, "hello.txt", filename)

	data, err := io.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", string(data))
}
