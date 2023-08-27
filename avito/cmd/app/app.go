package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/Lalipopp4/avito/internal/logger"
	"github.com/Lalipopp4/avito/internal/service"
)

// starting function
func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for {
			<-c
			close(c)
			logger.Logger.Log("Server is stopped")
			service.Stop()
			os.Exit(1)
		}

	}()
	service.Handle()
	logger.Logger.Log("Server is running...")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		logger.Logger.Log(err)
		return
	}
}
