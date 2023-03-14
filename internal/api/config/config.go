package config

import (
	"log"
	"os"
)

var Debug bool

func Debugln(v ...any) {
	if !Debug {
		return
	}
	log.Println(v...)
}

func Debugf(format string, a ...any) {
	if !Debug {
		return
	}
	log.Printf(format, a...)
}

var WorkingDir, _ = os.Getwd()
