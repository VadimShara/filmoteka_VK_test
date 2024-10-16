package filmoteka

import "github.com/urfave/cli/v2"

var cmdFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "server-host",
		Usage:   "server host",
		EnvVars: []string{"SERVER_HOST"},
		Value:   "localhost:8080",
	},
	&cli.StringFlag{
		Name:    "filmoteka-db-host",
		Usage:   "filmoteka db host",
		EnvVars: []string{"FILMOTEKA_DB_HOST"},
		Value:   "localhost:5432",
	},
	&cli.StringFlag{
		Name:    "filmoteka-db-user",
		Usage:   "filmoteka db user",
		EnvVars: []string{"FILMOTEKA_DB_USER"},
		Value:   "user",
	},
	&cli.StringFlag{
		Name:    "filmoteka-db-pass",
		Usage:   "filmoteka db password",
		EnvVars: []string{"FILMOTEKA_DB_PASS"},
		Value:   "pass",
	},
	&cli.StringFlag{
		Name:    "filmoteka-db-name",
		Usage:   "filmoteka db name",
		EnvVars: []string{"FILMOTEKA_DB_NAME"},
		Value:   "vk",
	},
	&cli.StringFlag{
		Name:    "filmoteka-db-sslmode",
		Usage:   "filmoteka db sslmode",
		EnvVars: []string{"FILMOTEKA_DB_SSLMODE"},
		Value:   "disable",
	},
}
