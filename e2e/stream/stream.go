package stream

import (
	"context"
	"fmt"
	"time"
)

type HttpEventServiceImpl struct {
}

var _ HttpEventService = (*HttpEventServiceImpl)(nil)

func (s *HttpEventServiceImpl) GetRandomValues(ctx context.Context) (values <-chan string, err error) {
	results := make(chan string, 10)

	go func() {
		defer close(results)
		count := 0

		for {
			select {
			case <-ctx.Done():
				return
			case results <- fmt.Sprintf("Hello %d", count):
				count++
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	return results, nil
}
