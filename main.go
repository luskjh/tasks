package main

import (
	"flag"
	"fmt"
	"os"

	"sorantics.com/tasks/internal/lib"
)

func serve(config *lib.Config) error {
	fmt.Printf("serving on port %d\n", config.Port)
	return nil
}

func run() error {
	var path, task string

	flag.StringVar(&path, "path", "config.json", "path to config file")
	flag.StringVar(&task, "task", "usage", "task to run")
	flag.Parse()

	config, err := lib.NewConfig(path)
	if err != nil {
		return err
	}

	switch task {
	case "serve":
		return serve(config)
	default:
		flag.Usage()
		return nil
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
