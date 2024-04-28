package download

import (
	"context"
	"io"
	"strings"
)

type HttpDownloadServiceImpl struct {
}

var _ HttpDownloadService = (*HttpDownloadServiceImpl)(nil)

func (s *HttpDownloadServiceImpl) Get(ctx context.Context) (asset io.Reader, filename string, contentType string, err error) {
	return strings.NewReader("Hello, World!"), "hello.txt", "text/plain", nil
}
