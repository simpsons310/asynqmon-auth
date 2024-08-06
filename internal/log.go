package asynqmonauth

import (
	"log"
	"os"
)

type Logger interface {
	Print(v ...any)
	Printf(format string, v ...any)
	Println(v ...any)
	Fatal(v ...any)
	Fatalf(format string, v ...any)
	Fatalln(v ...any)
	Panic(v ...any)
	Panicf(format string, v ...any)
	Panicln(v ...any)
}

func NewLogger() Logger {
	return log.New(os.Stdout, "[server] ", log.LstdFlags)
}
