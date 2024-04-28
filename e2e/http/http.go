package http

import (
	"context"
)

type HttpPeopleServiceImpl struct {
}

var _ HttpPeopleService = (*HttpPeopleServiceImpl)(nil)

func (s *HttpPeopleServiceImpl) GetRandom(ctx context.Context, age int8) (person *Person, err error) {
	if age < 0 {
		return nil, ErrAgen
	}

	return &Person{
		Name:    "Ella",
		Age:     age,
		Emotion: Emotion_Excited,
	}, nil
}
