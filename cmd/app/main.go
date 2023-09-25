package main

import (
	"log/slog"

	"github.com/Lalipopp4/test_api/internal/transport/rest"
)

func main() {
	server, err := rest.New()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if err := server.Run(); err != nil {
		slog.Error(err.Error())
	}
}
