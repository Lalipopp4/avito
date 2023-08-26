package logger

import (
	"log"
	"os"
)

type logger struct {
	log_file *os.File
}

var Logger = NewLogger(`C:\Users\anton\Go\src\github.com\Lalipopp4\avito\internal\logger\logfile.txt`)

func NewLogger(filename string) *logger {
	file, err := os.Open(filename)
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
	//log.Println(message...)
	log.SetOutput(l.log_file)
	log.Println(message...)

}

func (l logger) Fatal(message ...any) {
	log.SetOutput(os.Stdout)
	log.Println(message...)
	log.SetOutput(l.log_file)
	log.Println(message...)
}

func (l logger) Close() {
	l.log_file.Close()
}
