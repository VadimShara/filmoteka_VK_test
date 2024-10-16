package filmoteka

import (
	movie "vk-test-task/internal/store/movie"
	star "vk-test-task/internal/store/star"
)

type (
	Service interface {
		starsService
		moviesService
	}

	serviceImpl struct {
		starsStore  star.Store
		moviesStore movie.Store
	}
)

func New(
	stars star.Store,
	movies movie.Store,
) Service {
	return &serviceImpl{
		starsStore:  stars,
		moviesStore: movies,
	}
}
