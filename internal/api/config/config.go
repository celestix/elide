package config

import "log"

var Debug bool

func Debugln(v ...any) {
	log.Println(v...)
}

func Debugf(format string, a ...any) {
	log.Printf(format, a...)
}
