package http

import (
	"context"
	"net/http"
)

type HttpPeopleServiceImpl struct {
}

var _ HttpPeopleService = (*HttpPeopleServiceImpl)(nil)

func (s *HttpPeopleServiceImpl) GetRandom(ctx context.Context, age int8) (person *Person, err error) {
	if age < 0 {
		return nil, Errorf(http.StatusBadRequest, nil, "age must be greater than 0")
	}

	return &Person{
		Name:    "Ella",
		Age:     age,
		Emotion: Emotion_Excited,
	}, nil
}
