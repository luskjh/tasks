package main

import (
	"flag"
	"fmt"
)

func main() {
	var path, task string

	flag.StringVar(&path, "path", "config.json", "path to config file")
	flag.StringVar(&task, "task", "usage", "task to run")
	flag.Parse()

	fmt.Printf("path = %s, task = %s\n", path, task)
}
