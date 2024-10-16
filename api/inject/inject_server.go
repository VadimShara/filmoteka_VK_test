package inject

import (
	"vk-test-task/api/rest/handlers"
	"vk-test-task/internal/service/auth"
	"vk-test-task/internal/service/filmoteka"

	"github.com/google/wire"
	"github.com/urfave/cli/v2"
)

// wire set for loading the server.
var serverSet = wire.NewSet( // nolint
	provideResolver,
)

func provideResolver(c *cli.Context, filmotekaService filmoteka.Service, authService auth.Service) *handlers.Resolver {
	return handlers.NewResolver(c.String("server-host"), filmotekaService, authService)
}
