package upload

import (
	"context"
	"errors"
	"io"
)

type HttpStorageServiceImpl struct {
}

var _ HttpStorageService = (*HttpStorageServiceImpl)(nil)

func (s *HttpStorageServiceImpl) UploadFiles(ctx context.Context, files func() (string, io.Reader, error), id string) (results []*File, err error) {
	results = make([]*File, 0)

	for {
		filename, content, err := files()
		if errors.Is(err, io.EOF) {
			return results, nil
		} else if err != nil {
			return nil, err
		}

		size, _ := io.Copy(io.Discard, content)

		results = append(results, &File{
			Name: filename,
			Size: size,
		})
	}
}
