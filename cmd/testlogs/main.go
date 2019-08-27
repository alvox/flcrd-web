package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

func main() {
	path := flag.String("path", "/mnt/logs", "Application port")
	flag.Parse()

	dirs, err := filepath.Glob(*path)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dirs)
}
