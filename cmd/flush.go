package cmd

import (
	"elide/internal/api/config"
	"log"
	"os"
	"time"
)

func runFlush() {
	defer log.Println("[Elide][Flush]: Started")
	for {
		// clean in every 3? days
		time.Sleep(time.Hour * 24 * 3)
		dir := config.WorkingDir + "/downloads/photos"
		_ = os.RemoveAll(dir)
		_ = os.MkdirAll(dir, os.ModeDir)
	}
}
