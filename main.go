package main

import (
	"flag"
	"fmt"
	htemplate "html/template"
	"net/http"
	"os"

	"sorantics.com/tasks/internal/domain/dashboard"
	"sorantics.com/tasks/internal/lib"
	"sorantics.com/tasks/internal/static"
	"sorantics.com/tasks/internal/template"
)

func serve(config *lib.Config) error {
	distFS, err := static.DistFS(static.FS)
	if err != nil {
		return err
	}

	vite, err := lib.NewVite("internal/static/.vite/manifest.json")
	if err != nil {
		return err
	}

	funcMap := htemplate.FuncMap{"vite_entry": vite.Entry}
	reg, err := lib.NewRegistry(template.FS, funcMap)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", dashboard.Overview(reg))
	mux.Handle("GET /", http.FileServer(http.FS(distFS)))

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
