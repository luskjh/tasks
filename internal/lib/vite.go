package lib

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"html/template"
	"os"
)

//go:embed vite.html
var viteTemplateSource string
var viteTemplate = template.Must(template.New("vite.html").Parse(viteTemplateSource))

type Entry struct {
	File    string
	CSS     []string
	Imports []string
	IsEntry bool
}

type Vite map[string]template.HTML

func NewVite(path string) (*Vite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest map[string]*Entry
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	entries := make(Vite)
	for entry, item := range manifest {
		if item.IsEntry {
			entries[entry] = renderEntry(manifest, entry)
		}
	}

	return &entries, nil
}

func (v *Vite) Entry(name string) template.HTML {
	html, _ := (*v)[name]
	return html
}

func renderEntry(manifest map[string]*Entry, name string) template.HTML {
	seen := make(map[string]struct{})

	root := manifest[name]
	var modules, styles []string
	seen[name] = struct{}{}

	var walk func(imports []string)
	walk = func(imports []string) {
		for _, imported := range imports {
			_, importSeen := seen[imported]
			chunk, chunkExists := manifest[imported]

			if importSeen || !chunkExists {
				continue
			}
			seen[imported] = struct{}{}

			walk(chunk.Imports)
			modules = append(modules, chunk.File)
			styles = append(styles, chunk.CSS...)
		}
	}
	walk(root.Imports)

	styles = append(styles, root.CSS...)

	var buf bytes.Buffer
	_ = viteTemplate.Execute(&buf, struct {
		Script  string
		Modules []string
		Styles  []string
	}{
		Script:  root.File,
		Modules: modules,
		Styles:  styles,
	})

	return template.HTML(buf.String())
}
