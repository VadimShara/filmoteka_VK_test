package main

import (
	"log"
	"os"

	"vk-test-task/cmd/filmoteka"
	"vk-test-task/pkg/logger"

	"github.com/urfave/cli/v2"
)

var globalFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "log-level",
		Usage:   "log level",
		EnvVars: []string{"LOG_LEVEL"},
		Value:   "local",
	},
}

// @Version 1.0.0
// @Title VK Filmoteka API Service
// @Description VK Filmoteka API Service Documentation
// @Server http://localhost:8080 local server
// @Security AuthorizationHeader bearer token
// @SecurityScheme AuthorizationHeader http bearer JWT-token
func main() {
	app := &cli.App{
		Usage: "Filmoteka API Service",
		Commands: []*cli.Command{
			&filmoteka.Cmd,
		},
		Flags: globalFlags,
		Before: cli.BeforeFunc(func(c *cli.Context) error {
			return logger.SetupLogger(c.String("log-level"))
		}),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("error start service: %v", err)
	}
}
