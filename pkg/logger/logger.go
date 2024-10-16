package logger

import (
	"fmt"
	"os"

	"golang.org/x/exp/slog"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

var Log *slog.Logger

func SetupLogger(env string) error {
	switch env {
	case envLocal:
		Log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		Log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		return fmt.Errorf("invalid logger level: %s", env)
	}
	return nil
}
