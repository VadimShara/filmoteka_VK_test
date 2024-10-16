package inject

import (
	"fmt"

	"vk-test-task/internal/store/movie"
	"vk-test-task/internal/store/star"
	"vk-test-task/internal/store/user"
	"vk-test-task/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/urfave/cli/v2"

	"github.com/google/wire"
)

var storeSet = wire.NewSet( // nolint
	createDBClient,
	provideStores,
)

type stores struct {
	stars  star.Store
	movies movie.Store
	users  user.Store
}

func provideStores(c *cli.Context, db *pgxpool.Pool) stores {
	return stores{
		stars:  star.New(db),
		movies: movie.New(db),
		users:  user.New(db),
	}
}

func createDBClient(c *cli.Context) (*pgxpool.Pool, error) {
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		c.String("filmoteka-db-user"),
		c.String("filmoteka-db-pass"),
		c.String("filmoteka-db-host"),
		c.String("filmoteka-db-name"),
		c.String("filmoteka-db-sslmode"),
	)
	logger.Log.Debug("connecting to database", "url", databaseURL)

	postgresClient, err := pgxpool.New(c.Context, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := postgresClient.Ping(c.Context); err != nil {
		return nil, err
	}

	return postgresClient, nil
}
