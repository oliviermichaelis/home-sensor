package infrastructure

import "log"

type Logger struct{}

func (logger Logger) Log(args ...interface{}) {
	log.Println(args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	log.Fatal(args...)
}
