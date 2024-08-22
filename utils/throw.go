package utils

import (
	"log/slog"
	"os"
)

func Throw(err string) {
	slog.Error(err)
	os.Exit(1)
}
