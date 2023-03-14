package main

import (
	"elide/cmd"
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.MkdirAll("downloads/photos", os.ModeDir))
	cmd.Execute()
}
