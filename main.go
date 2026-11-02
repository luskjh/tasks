package main

import (
	"flag"
	"fmt"
	"os"

	"sorantics.com/tasks/internal/lib"
)

func run() error {
	var path, task string

	flag.StringVar(&path, "path", "config.json", "path to config file")
	flag.StringVar(&task, "task", "usage", "task to run")
	flag.Parse()

	config, err := lib.NewConfig(path)
	if err != nil {
		return err
	}

	fmt.Printf("config.port = %d\n", config.Port)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
