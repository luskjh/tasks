package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"sorantics.com/tasks/internal/domain/dashboard"
	"sorantics.com/tasks/internal/lib"
	"sorantics.com/tasks/internal/template"
)

func serve(config *lib.Config) error {
	reg, err := lib.NewRegistry(template.FS)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", dashboard.Overview(reg))

	addr := config.Addr()
	fmt.Printf("listening on %s\n", addr)
	return http.ListenAndServe(addr, mux)
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
