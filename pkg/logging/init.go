package logging

import (
	"log/slog"
	"os"
)

const (
	PATHLOGFILE = "pkg/file/log/log.txt"
)

type logger struct {
	logStd  *slog.Logger
	logFile *slog.Logger
}

func New() (Logger, error) {
	file, err := os.OpenFile(PATHLOGFILE, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	return &logger{
		slog.New(slog.NewTextHandler(os.Stdout, nil)),
		slog.New(slog.NewTextHandler(file, nil)),
	}, nil

}
