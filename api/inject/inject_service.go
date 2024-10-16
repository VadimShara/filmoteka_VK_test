package inject

import (
	"vk-test-task/internal/service/auth"
	"vk-test-task/internal/service/filmoteka"

	"github.com/google/wire"
)

var serviceSet = wire.NewSet( // nolint
	provideFilmotekaService,
	provideAuthService,
)

func provideFilmotekaService(s stores) filmoteka.Service {
	return filmoteka.New(s.stars, s.movies)
}

func provideAuthService(s stores) (auth.Service, error) {
	return auth.New(s.users)
}
