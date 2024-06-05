package logger

import (
	"log/slog"

	"github.com/oke11o/wslog"
)

func New(asJson bool, level slog.Leveler) *slog.Logger {
	return wslog.New(asJson, level)
}
