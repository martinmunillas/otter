package log

import (
	"log/slog"
	"os"

	"github.com/a-h/templ/cmd/templ/sloghandler"
)

func NewLogger(verbose bool) *slog.Logger {
	loggingLevel := slog.LevelInfo.Level()
	if verbose {
		loggingLevel = slog.LevelDebug.Level()
	}
	logger := slog.New(sloghandler.NewHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: verbose,
		Level:     loggingLevel,
	}))
	return logger
}
