package logger

import (
	"log"
	"os"
)

// implementation of logger to make logging in os.Stdin
// and logging file easier

type logger struct {
	LogFile *os.File
}

var Logger = NewLogger("logger/logfile.txt")

func NewLogger(filename string) *logger {
	file, err := os.OpenFile("internal/logger/logfile.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &logger{
		file,
	}
}

func (l logger) Log(message ...any) {
	log.SetOutput(os.Stdout)
	log.Println(message...)
	log.SetOutput(l.LogFile)
	log.Println(message...)

}

func (l *logger) Close() {
	l.LogFile.Close()
}
