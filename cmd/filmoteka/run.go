package filmoteka

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"vk-test-task/api/inject"
	"vk-test-task/pkg/logger"

	"github.com/urfave/cli/v2"
)

var Cmd = cli.Command{
	Name:   "filmoteka",
	Usage:  "Run filmoteka API",
	Flags:  cmdFlags,
	Action: run,
}

func run(c *cli.Context) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
			return
		case s := <-sig:
			logger.Log.Info("received", "signal", s.String())
			cancel()
		}
	}()

	app, err := inject.InitializeApplication(c, ctx)
	if err != nil {
		logger.Log.Error("main: cannot initialize server", "error", err.Error())
		os.Exit(1)
	}
	logger.Log.Info("server started", "address", app.Resolver.GetAddr())

	app.Resolver.Run() // nolint

	<-ctx.Done()
	logger.Log.Debug("ctx end received")

	return nil
}
