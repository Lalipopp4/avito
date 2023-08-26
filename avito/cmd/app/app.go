package main

import (
	"log"

	"net/http"

	//"github.com/Lalipopp4/avito/internal/logger"
	"github.com/Lalipopp4/avito/internal/service"
)

func main() {
	//logger := logger.NewLogger("github.com/Lalipopp4/avito/internal/logger/logfile.txt")
	service.Handle()
	//logger.Logger.Log("Server is running...")
	log.Println("Server is running...")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		//logger.Logger.Fatal(err)
	}
	//logger.Logger.Log("Server is stopping...")
	service.Stop()

}
